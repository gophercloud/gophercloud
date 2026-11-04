// Package impl extracts a model.API from Gophercloud Go source via static
// analysis. It discovers operations (HTTP method + URL), request/response
// fields and query parameters, including the irregular cases that runtime
// reflection cannot see: manually-assembled map bodies, json:"-" fields, and
// URLs composed through helper functions.
package impl

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/gophercloud/gophercloud/tools/apicheck/model"
)

// verbMethods maps ServiceClient method names to HTTP verbs.
var verbMethods = map[string]string{
	"Get": "GET", "Post": "POST", "Put": "PUT",
	"Patch": "PATCH", "Delete": "DELETE", "Head": "HEAD",
}

// bodyArgIndex is the positional index (after ctx, url) of the request body for
// methods that carry one.
var bodyArgIndex = map[string]int{"POST": 2, "PUT": 2, "PATCH": 2}

// Unresolved records an operation whose URL could not be statically resolved.
type Unresolved struct {
	Package  string
	Function string
	Reason   string
}

// Result is the outcome of analysing one service.
type Result struct {
	API        *model.API
	Unresolved []Unresolved
}

// Load analyses every package under gophercloudRoot/pkgPath and returns the
// discovered API for the given service.
func Load(service, gophercloudRoot, pkgPath string) (*Result, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax |
			packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps | packages.NeedFiles,
		Dir: gophercloudRoot,
	}
	pkgs, err := packages.Load(cfg, "./"+pkgPath+"/...")
	if err != nil {
		return nil, fmt.Errorf("loading packages under %s: %w", pkgPath, err)
	}

	res := &Result{API: &model.API{Service: service}}
	for _, pkg := range pkgs {
		if len(pkg.Syntax) == 0 {
			continue
		}
		a := &analyzer{service: service, pkg: pkg}
		a.prepare()
		a.analyze(res)
	}

	sort.Slice(res.API.Operations, func(i, j int) bool {
		return res.API.Operations[i].Key() < res.API.Operations[j].Key()
	})
	return res, nil
}

type analyzer struct {
	service string
	pkg     *packages.Package

	funcs      map[string]*ast.FuncDecl // top-level funcs by name
	methods    map[string]*types.Named  // To*Map/To*Query method -> owner type
	methodDecl map[string]*ast.FuncDecl // To* method name -> decl
	urlCache   map[string][]string      // url helper func name -> segments
	respPool   []model.Field            // union of response struct fields
	strVars    map[string]string        // package string var/const -> value
}

func (a *analyzer) prepare() {
	a.funcs = map[string]*ast.FuncDecl{}
	a.methods = map[string]*types.Named{}
	a.methodDecl = map[string]*ast.FuncDecl{}
	a.urlCache = map[string][]string{}
	a.strVars = map[string]string{}

	for _, f := range a.pkg.Syntax {
		for _, decl := range f.Decls {
			if gd, ok := decl.(*ast.GenDecl); ok && (gd.Tok == token.VAR || gd.Tok == token.CONST) {
				a.collectStringVars(gd)
				continue
			}
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fd.Recv == nil {
				a.funcs[fd.Name.Name] = fd
				continue
			}
			// Method: record To*Map / To*Query builders by name.
			name := fd.Name.Name
			if strings.HasPrefix(name, "To") && (strings.HasSuffix(name, "Map") || strings.HasSuffix(name, "Query")) {
				if named := a.recvNamed(fd); named != nil {
					a.methods[name] = named
					a.methodDecl[name] = fd
				}
			}
		}
	}
	a.buildResponsePool()
}

// buildResponsePool collects the union of json-tagged fields across every
// exported struct in the package that is not a request-options type. This is a
// proxy for "does Gophercloud model this response field anywhere in the
// package"; it deliberately over-approximates to avoid false negatives, since
// impl response structs cannot be reliably tied to individual operations.
func (a *analyzer) buildResponsePool() {
	scope := a.pkg.Types.Scope()
	seen := map[string]bool{}
	for _, name := range scope.Names() {
		if strings.HasSuffix(name, "Opts") || strings.HasSuffix(name, "OptsBuilder") {
			continue
		}
		tn, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}
		named, ok := tn.Type().(*types.Named)
		if !ok {
			continue
		}
		if _, ok := named.Underlying().(*types.Struct); !ok {
			continue
		}
		for _, f := range a.structFields(named, "json") {
			if seen[f.Name] {
				continue
			}
			seen[f.Name] = true
			a.respPool = append(a.respPool, f)
		}
	}
}

func (a *analyzer) recvNamed(fd *ast.FuncDecl) *types.Named {
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return nil
	}
	t := a.pkg.TypesInfo.TypeOf(fd.Recv.List[0].Type)
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if named, ok := t.(*types.Named); ok {
		return named
	}
	return nil
}

// analyze walks every exported top-level function, emitting an Operation for
// each ServiceClient verb call it makes.
func (a *analyzer) analyze(res *Result) {
	for name, fd := range a.funcs {
		if !ast.IsExported(name) || fd.Body == nil {
			continue
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			// List operations construct a pagination.Pager rather than calling
			// a verb method directly: pagination.NewPager(client, url, ...).
			if sel.Sel.Name == "NewPager" && len(call.Args) >= 2 {
				a.emit(res, name, fd, "GET", call.Args[1], nil)
				return true
			}
			// Generic verb call: client.Request(ctx, "PUT", url, opts). The verb
			// is the first string argument and the URL the second. The body is
			// carried inside RequestOpts rather than as a positional argument, so
			// it is not extracted here.
			if sel.Sel.Name == "Request" && a.isServiceClient(sel.X) && len(call.Args) >= 3 {
				if verb, ok := a.constString(call.Args[1]); ok {
					a.emit(res, name, fd, strings.ToUpper(verb), call.Args[2], nil)
				}
				return true
			}
			verb, ok := verbMethods[sel.Sel.Name]
			if !ok || !a.isServiceClient(sel.X) {
				return true
			}
			var body ast.Expr
			if idx, has := bodyArgIndex[verb]; has && len(call.Args) > idx {
				body = call.Args[idx]
			}
			a.emit(res, name, fd, verb, call.Args[1], body)
			return true
		})
	}
}

// emit builds an Operation from a verb call or NewPager call. urlExpr is the
// URL argument; body is the request body expression (nil if none).
func (a *analyzer) emit(res *Result, fnName string, fd *ast.FuncDecl, verb string, urlExpr, body ast.Expr) {
	segs, ok := a.urlFromExpr(fd, urlExpr)
	if !ok {
		res.Unresolved = append(res.Unresolved, Unresolved{
			Package: a.pkg.PkgPath, Function: fnName, Reason: "unresolved URL expression",
		})
		return
	}
	p := normalizePath("/" + strings.Join(segs, "/"))

	op := model.Operation{
		Service: a.service,
		Method:  verb,
		Path:    p,
		Source:  a.pkg.Name + "." + fnName,
	}

	if body != nil {
		action, fields := a.bodyInfo(fd, body)
		op.RequestFields = fields
		if strings.HasSuffix(p, "/action") {
			op.Action = action
		}
	}

	// Query params: look for a To*Query builder call in the function body.
	op.QueryParams = a.queryFor(fd)

	// Response fields: the package-wide pool (see buildResponsePool).
	op.ResponseFields = a.respPool

	res.API.Operations = append(res.API.Operations, op)
}

// urlFromExpr resolves a URL argument, following local variable assignments
// (e.g. "url := listURL(client)") to the underlying helper call.
func (a *analyzer) urlFromExpr(fd *ast.FuncDecl, expr ast.Expr) ([]string, bool) {
	if id, ok := expr.(*ast.Ident); ok {
		if rhs := a.localAssign(fd, id.Name); rhs != nil {
			return a.urlFromExpr(fd, rhs)
		}
		return nil, false
	}
	return a.resolveURL(expr)
}

// localAssign returns the RHS of the first "name := rhs" or "name = rhs" in fd.
func (a *analyzer) localAssign(fd *ast.FuncDecl, name string) ast.Expr {
	var rhs ast.Expr
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if rhs != nil {
			return false
		}
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for i, lhs := range as.Lhs {
			id, ok := lhs.(*ast.Ident)
			if ok && id.Name == name && i < len(as.Rhs) {
				rhs = as.Rhs[i]
				return false
			}
		}
		return true
	})
	return rhs
}

// isServiceClient reports whether expr has type *gophercloud.ServiceClient.
func (a *analyzer) isServiceClient(expr ast.Expr) bool {
	t := a.pkg.TypesInfo.TypeOf(expr)
	if t == nil {
		// Fallback to the naming convention.
		id, ok := expr.(*ast.Ident)
		return ok && id.Name == "client"
	}
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return named.Obj().Name() == "ServiceClient" &&
		strings.Contains(named.Obj().Pkg().Path(), "gophercloud")
}

// resolveURL turns a URL argument expression into path segments, with path
// parameters normalised to "{}".
func (a *analyzer) resolveURL(expr ast.Expr) ([]string, bool) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "ServiceURL" {
			return a.argsToSegments(e.Args), true
		}
		if id, ok := e.Fun.(*ast.Ident); ok {
			return a.resolveURLFunc(id.Name)
		}
	case *ast.SelectorExpr:
		// client.Endpoint is the service root (e.g. Swift's account endpoint,
		// which already embeds /v1/{account}): zero path segments.
		if e.Sel.Name == "Endpoint" && a.isServiceClient(e.X) {
			return []string{}, true
		}
	case *ast.BinaryExpr:
		// URL built by concatenation, e.g. client.Endpoint + "?bulk-delete=true".
		// The right operand is a query-string/suffix that does not affect the
		// path, so resolve the left operand alone.
		if segs, ok := a.resolveURL(e.X); ok {
			return segs, true
		}
	}
	return nil, false
}

// resolveURLFunc resolves a local URL-builder helper to its segment template.
func (a *analyzer) resolveURLFunc(name string) ([]string, bool) {
	if segs, ok := a.urlCache[name]; ok {
		return segs, segs != nil
	}
	a.urlCache[name] = nil // guard against cycles
	fd, ok := a.funcs[name]
	if !ok || fd.Body == nil {
		return nil, false
	}
	ret := lastReturn(fd.Body)
	if ret == nil || len(ret.Results) == 0 {
		return nil, false
	}
	segs, ok := a.resolveURL(ret.Results[0])
	if ok {
		a.urlCache[name] = segs
	}
	return segs, ok
}

// argsToSegments converts ServiceURL arguments to segments: constant strings
// (literals and named string constants) become their value, everything else
// becomes the "{}" wildcard.
func (a *analyzer) argsToSegments(args []ast.Expr) []string {
	var segs []string
	for _, arg := range args {
		if v, ok := a.constString(arg); ok {
			segs = append(segs, v)
			continue
		}
		segs = append(segs, "{}")
	}
	return segs
}

// constString resolves an expression to its constant string value if possible.
func (a *analyzer) constString(expr ast.Expr) (string, bool) {
	if tv, ok := a.pkg.TypesInfo.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		return constant.StringVal(tv.Value), true
	}
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if v, err := strconv.Unquote(lit.Value); err == nil {
			return v, true
		}
	}
	// Package-level string var/const not folded by the type checker.
	if id, ok := expr.(*ast.Ident); ok {
		if v, ok := a.strVars[id.Name]; ok {
			return v, true
		}
	}
	return "", false
}

// collectStringVars records package-level "name = \"literal\"" declarations.
func (a *analyzer) collectStringVars(gd *ast.GenDecl) {
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, name := range vs.Names {
			if i >= len(vs.Values) {
				continue
			}
			if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if v, err := strconv.Unquote(lit.Value); err == nil {
					a.strVars[name.Name] = v
				}
			}
		}
	}
}

// bodyInfo extracts the action name and request fields from a body argument.
func (a *analyzer) bodyInfo(fd *ast.FuncDecl, body ast.Expr) (string, []model.Field) {
	switch e := body.(type) {
	case *ast.CompositeLit:
		// Inline map literal, e.g. map[string]any{"forceDelete": ""}.
		return mapLitInfo(e)
	case *ast.Ident:
		if e.Name == "nil" {
			return "", nil
		}
		// Body may be a locally-assigned map literal, e.g.
		//   b := map[string]any{"resetNetwork": nil}
		if cl := a.localMapLit(fd, e.Name); cl != nil {
			return mapLitInfo(cl)
		}
		// Otherwise it is typically assigned from opts.To*Map(). Resolve the
		// builder from the variable's own defining assignment so we pick the
		// correct one when a function calls several To*Map builders (e.g.
		// servers.Create also merges in ToSchedulerHintsMap).
		if rhs := a.localAssign(fd, e.Name); rhs != nil {
			if mname := builderCallName(rhs); mname != "" {
				if _, known := a.methods[mname]; known {
					action, fields := a.builderInfo(mname)
					// Additional builders whose output is merged into the body
					// via maps.Copy(body, x) contribute more fields (e.g.
					// scheduler hints copied into the server create body).
					for _, extra := range a.mergedBuilders(fd, e.Name) {
						if _, xfields := a.builderInfo(extra); len(xfields) > 0 {
							fields = append(fields, xfields...)
						}
					}
					return action, dedupeFields(fields)
				}
			}
		}
		// Fallback: any To*Map builder in the function body.
		if mname := a.findBuilderCall(fd, "Map"); mname != "" {
			return a.builderInfo(mname)
		}
	}
	return "", nil
}

// mergedBuilders finds To*Map builders whose result is merged into the body
// variable via maps.Copy(body, x), and returns their method names. This captures
// secondary request bodies such as scheduler hints copied into a server create
// body.
func (a *analyzer) mergedBuilders(fd *ast.FuncDecl, body string) []string {
	var out []string
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Copy" || len(call.Args) != 2 {
			return true
		}
		// maps.Copy(dst, src): dst must be the body variable.
		dst, ok := call.Args[0].(*ast.Ident)
		if !ok || dst.Name != body {
			return true
		}
		src, ok := call.Args[1].(*ast.Ident)
		if !ok {
			return true
		}
		if rhs := a.localAssign(fd, src.Name); rhs != nil {
			if mname := builderCallName(rhs); mname != "" {
				if _, known := a.methods[mname]; known {
					out = append(out, mname)
				}
			}
		}
		return true
	})
	return out
}

// localMapLit finds a local assignment "name := map[...]{...}" within fd.
func (a *analyzer) localMapLit(fd *ast.FuncDecl, name string) *ast.CompositeLit {
	var found *ast.CompositeLit
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for i, lhs := range as.Lhs {
			id, ok := lhs.(*ast.Ident)
			if !ok || id.Name != name || i >= len(as.Rhs) {
				continue
			}
			if cl, ok := as.Rhs[i].(*ast.CompositeLit); ok {
				if _, isMap := cl.Type.(*ast.MapType); isMap {
					found = cl
				}
			}
		}
		return true
	})
	return found
}

// mapLitInfo returns the single action key and any string keys of a map literal.
func mapLitInfo(cl *ast.CompositeLit) (string, []model.Field) {
	var keys []string
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if lit, ok := kv.Key.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if v, err := strconv.Unquote(lit.Value); err == nil {
				keys = append(keys, v)
			}
		}
	}
	action := ""
	if len(keys) == 1 {
		action = keys[0]
	}
	// A single-key wrapper contributes no directly-visible fields.
	var fields []model.Field
	if len(keys) > 1 {
		for _, k := range keys {
			fields = append(fields, model.Field{Name: k})
		}
	}
	return action, fields
}

// builderCallName returns the method name of an expression of the form
// x.To<X>Map() / x.To<X>Query(), or "" if expr is not such a call.
func builderCallName(expr ast.Expr) string {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return ""
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	name := sel.Sel.Name
	if strings.HasPrefix(name, "To") && (strings.HasSuffix(name, "Map") || strings.HasSuffix(name, "Query")) {
		return name
	}
	return ""
}

// findBuilderCall finds the first call like opts.To<X><suffix>() within fd and
// returns the method name.
func (a *analyzer) findBuilderCall(fd *ast.FuncDecl, suffix string) string {
	var found string
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if strings.HasPrefix(sel.Sel.Name, "To") && strings.HasSuffix(sel.Sel.Name, suffix) {
			if _, known := a.methods[sel.Sel.Name]; known {
				found = sel.Sel.Name
				return false
			}
		}
		return true
	})
	return found
}

// builderInfo resolves the action wrapping key and request fields for a
// To*Map builder method.
func (a *analyzer) builderInfo(methodName string) (string, []model.Field) {
	named := a.methods[methodName]
	decl := a.methodDecl[methodName]
	parent := ""
	var manualKeys []string
	if decl != nil {
		parent, manualKeys = a.builderBody(decl)
	}
	fields := a.structFields(named, "json")
	for _, k := range manualKeys {
		fields = append(fields, model.Field{Name: k})
	}
	return parent, dedupeFields(fields)
}

// builderBody scans a builder method for the BuildRequestBody parent key and
// any manual b["key"] = ... assignments.
func (a *analyzer) builderBody(decl *ast.FuncDecl) (parent string, manualKeys []string) {
	if decl.Body == nil {
		return "", nil
	}
	ast.Inspect(decl.Body, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.CallExpr:
			if sel, ok := e.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "BuildRequestBody" && len(e.Args) >= 2 {
				if v, ok := a.constString(e.Args[1]); ok && v != "" {
					parent = v
				}
			}
		case *ast.ReturnStmt:
			// return map[string]any{"rebuild": b} — the single key is the
			// wrapping/action name when BuildRequestBody used an empty parent.
			if parent == "" && len(e.Results) > 0 {
				if cl, ok := e.Results[0].(*ast.CompositeLit); ok {
					if _, isMap := cl.Type.(*ast.MapType); isMap {
						if k, _ := mapLitInfo(cl); k != "" {
							parent = k
						}
					}
				}
			}
		case *ast.AssignStmt:
			// b["key"] = ...
			for _, lhs := range e.Lhs {
				if ix, ok := lhs.(*ast.IndexExpr); ok {
					if lit, ok := ix.Index.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						if v, err := strconv.Unquote(lit.Value); err == nil {
							manualKeys = append(manualKeys, v)
						}
					}
				}
			}
		}
		return true
	})
	return parent, manualKeys
}

// queryFor finds a To*Query builder in fd and returns its query params.
func (a *analyzer) queryFor(fd *ast.FuncDecl) []model.Field {
	mname := a.findBuilderCall(fd, "Query")
	if mname == "" {
		return nil
	}
	return a.structFields(a.methods[mname], "q")
}

// structFields returns the exported fields of a named struct type, keyed by the
// given tag ("json" or "q"). Embedded structs are flattened. Fields tagged "-"
// are recorded as ManualHandled (json) or skipped (q).
func (a *analyzer) structFields(named *types.Named, tag string) []model.Field {
	if named == nil {
		return nil
	}
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	var out []model.Field
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Exported() {
			continue
		}
		tags := reflect.StructTag(st.Tag(i))
		if f.Embedded() {
			if en, ok := f.Type().(*types.Named); ok {
				out = append(out, a.structFields(en, tag)...)
			}
			continue
		}
		val := tags.Get(tag)
		name := strings.Split(val, ",")[0]
		if name == "-" {
			if tag == "json" {
				out = append(out, model.Field{Name: f.Name(), ManualHandled: true})
			}
			continue
		}
		if name == "" {
			// No tag: query params require one; json defaults to field name.
			if tag == "q" {
				continue
			}
			name = f.Name()
		}
		field := model.Field{Name: name, Required: tags.Get("required") == "true"}
		out = append(out, field)
	}
	return out
}

func dedupeFields(fields []model.Field) []model.Field {
	seen := map[string]bool{}
	var out []model.Field
	for _, f := range fields {
		if seen[f.Name] {
			continue
		}
		seen[f.Name] = true
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func lastReturn(body *ast.BlockStmt) *ast.ReturnStmt {
	var ret *ast.ReturnStmt
	for _, stmt := range body.List {
		if r, ok := stmt.(*ast.ReturnStmt); ok {
			ret = r
		}
	}
	return ret
}

// normalizePath collapses duplicate slashes for consistent matching.
func normalizePath(p string) string {
	return path.Clean(p)
}
