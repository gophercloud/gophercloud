// Package spec builds a model.API from an OpenStack OpenAPI 3.1 document.
//
// It flattens microversion "oneOf" unions into a single set of fields, tagging
// each field/operation with the microversion (x-openstack.min-ver, or the
// "_<mv>" suffix on component schema names) that introduced it.
package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"go.yaml.in/yaml/v4"

	"github.com/gophercloud/gophercloud/tools/apicheck/model"
)

// maxDepth bounds recursion when flattening nested schemas.
const maxDepth = 6

// xOpenStack captures the vendor extension used throughout the OpenStack specs.
type xOpenStack struct {
	MinVer       string `yaml:"min-ver"`
	MaxVer       string `yaml:"max-ver"`
	ActionName   string `yaml:"action-name"`
	Discriminate string `yaml:"discriminator"`
}

// Load parses the OpenAPI document for a service and returns its model.API.
// service is the logical name, specFile is the path to the YAML document, and
// pathPrefix (e.g. "/v2.1") is stripped from every path. When stripProjectID is
// true, a leading project-scope path parameter left after the prefix strip is
// also removed and the resulting duplicate operations are merged.
func Load(service, specFile, pathPrefix string, stripProjectID bool) (*model.API, error) {
	data, err := os.ReadFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("reading spec %s: %w", specFile, err)
	}

	doc, err := libopenapi.NewDocument(data)
	if err != nil {
		return nil, fmt.Errorf("loading spec %s: %w", specFile, err)
	}
	docModel, err := doc.BuildV3Model()
	if err != nil {
		// Local-only specs still build; report but continue if a model exists.
		if docModel == nil {
			return nil, fmt.Errorf("building model for %s: %w", specFile, err)
		}
	}

	api := &model.API{Service: service}
	if docModel.Model.Info != nil {
		api.Version = docModel.Model.Info.Version
	}
	if docModel.Model.Paths == nil {
		return api, nil
	}

	e := &extractor{service: service}
	for path, item := range docModel.Model.Paths.PathItems.FromOldest() {
		rel := stripPrefix(path, pathPrefix)
		if stripProjectID {
			rel = stripLeadingParam(rel)
		}
		for method, op := range verbs(item) {
			if op == nil {
				continue
			}
			e.extract(api, method, rel, op)
		}
	}

	if stripProjectID {
		api.Operations = dedupeOps(api.Operations)
	}
	sort.Slice(api.Operations, func(i, j int) bool {
		return api.Operations[i].Key() < api.Operations[j].Key()
	})
	return api, nil
}

// LoadFromRoot is a convenience wrapper resolving specFile against a root dir.
func LoadFromRoot(service, specRoot, specFile, pathPrefix string, stripProjectID bool) (*model.API, error) {
	return Load(service, filepath.Join(specRoot, specFile), pathPrefix, stripProjectID)
}

type extractor struct {
	service string
}

// extract turns one OpenAPI operation into one or more model.Operations. A
// POST to a path ending in "/action" is expanded into one operation per action
// variant found in the request body's oneOf.
func (e *extractor) extract(api *model.API, method, path string, op *v3.Operation) {
	xo := readX(op.Extensions)
	query := e.queryParams(op)

	if method == "POST" && strings.HasSuffix(path, "/action") {
		byName := map[string]model.Operation{}
		var order []string
		for _, actOp := range e.expandActions(path, op, xo, query) {
			if prev, ok := byName[actOp.Action]; ok {
				// Same action at a newer microversion: keep the lowest
				// min-ver and merge the union of request fields.
				if lessVer(prev.MinVer, actOp.MinVer) {
					actOp.MinVer = prev.MinVer
				}
				actOp.RequestFields = dedupe(append(prev.RequestFields, actOp.RequestFields...))
			} else {
				order = append(order, actOp.Action)
			}
			byName[actOp.Action] = actOp
		}
		for _, name := range order {
			api.Operations = append(api.Operations, byName[name])
		}
		return
	}

	mop := model.Operation{
		Service:        e.service,
		Method:         method,
		Path:           path,
		MinVer:         xo.MinVer,
		OperationID:    op.OperationId,
		QueryParams:    query,
		RequestFields:  e.bodyFields(op.RequestBody),
		ResponseFields: e.responseFields(op),
	}
	api.Operations = append(api.Operations, mop)
}

// expandActions produces one Operation per action variant of an /action POST.
func (e *extractor) expandActions(path string, op *v3.Operation, xo xOpenStack, query []model.Field) []model.Operation {
	var out []model.Operation
	members := requestOneOf(op.RequestBody)
	if len(members) == 0 {
		// Single, non-union action body.
		if name, fields := actionFrom(nil); name != "" {
			out = append(out, e.actionOp(path, name, "", fields, query, op))
		}
		return out
	}
	for _, m := range members {
		name, fields := actionFrom(m)
		if name == "" {
			continue
		}
		mv := refMinVer(m)
		if x := readX(schemaExt(m)); x.MinVer != "" {
			mv = x.MinVer
		}
		out = append(out, e.actionOp(path, name, mv, fields, query, op))
	}
	return out
}

func (e *extractor) actionOp(path, action, minVer string, req, query []model.Field, op *v3.Operation) model.Operation {
	return model.Operation{
		Service:        e.service,
		Method:         "POST",
		Path:           path,
		Action:         action,
		MinVer:         minVer,
		OperationID:    op.OperationId,
		RequestFields:  req,
		QueryParams:    query,
		ResponseFields: e.responseFields(op),
	}
}

// actionFrom returns the action name (from x-openstack.action-name, else the
// single required/top-level property) and the action's request fields.
func actionFrom(m *base.SchemaProxy) (string, []model.Field) {
	if m == nil {
		return "", nil
	}
	s := m.Schema()
	if s == nil {
		return "", nil
	}
	if x := readX(s.Extensions); x.ActionName != "" {
		var fields []model.Field
		flatten(m, "", "", 0, map[string]bool{}, &fields)
		// Drop the wrapper property itself if it equals the action name.
		return x.ActionName, subFields(fields, x.ActionName)
	}
	// Fallback: single top-level property is the action.
	if s.Properties != nil && s.Properties.Len() == 1 {
		for name := range s.Properties.FromOldest() {
			var fields []model.Field
			flatten(m, "", "", 0, map[string]bool{}, &fields)
			return name, subFields(fields, name)
		}
	}
	return "", nil
}

// subFields returns the flattened fields whose JSONPath is nested under wrapper.
func subFields(all []model.Field, wrapper string) []model.Field {
	var out []model.Field
	for _, f := range all {
		if f.Name == wrapper {
			continue
		}
		out = append(out, f)
	}
	return out
}

func (e *extractor) queryParams(op *v3.Operation) []model.Field {
	var out []model.Field
	for _, p := range op.Parameters {
		if p == nil || p.In != "query" {
			continue
		}
		x := readX(p.Extensions)
		f := model.Field{Name: p.Name, MinVer: x.MinVer, MaxVer: x.MaxVer}
		if p.Required != nil {
			f.Required = *p.Required
		}
		if p.Deprecated {
			f.Deprecated = true
		}
		out = append(out, f)
	}
	return dedupe(out)
}

func (e *extractor) bodyFields(rb *v3.RequestBody) []model.Field {
	if rb == nil || rb.Content == nil {
		return nil
	}
	mt := jsonMedia(rb.Content)
	if mt == nil || mt.Schema == nil {
		return nil
	}
	var out []model.Field
	flatten(mt.Schema, "", "", 0, map[string]bool{}, &out)
	return dedupe(out)
}

// responseFields flattens the first successful (2xx) JSON response body.
func (e *extractor) responseFields(op *v3.Operation) []model.Field {
	if op.Responses == nil || op.Responses.Codes == nil {
		return nil
	}
	var resp *v3.Response
	for code, r := range op.Responses.Codes.FromOldest() {
		if strings.HasPrefix(code, "2") {
			resp = r
			break
		}
	}
	if resp == nil || resp.Content == nil {
		return nil
	}
	mt := jsonMedia(resp.Content)
	if mt == nil || mt.Schema == nil {
		return nil
	}
	var out []model.Field
	flatten(mt.Schema, "", "", 0, map[string]bool{}, &out)
	return dedupe(out)
}

// flatten walks a schema proxy, appending leaf fields to out. Wrapper objects
// (properties that are themselves objects with defined properties) are
// descended rather than emitted, so nested request/response envelopes collapse
// to their leaf property names. minVer is propagated from oneOf ref suffixes.
func flatten(proxy *base.SchemaProxy, prefix, minVer string, depth int, visited map[string]bool, out *[]model.Field) {
	if proxy == nil || depth > maxDepth {
		return
	}
	if ref := proxy.GetReference(); ref != "" {
		if visited[ref] {
			return
		}
		visited[ref] = true
		if mv := mvFromRef(ref); mv != "" && minVer == "" {
			minVer = mv
		}
	}
	s := proxy.Schema()
	if s == nil {
		return
	}
	if x := readX(s.Extensions); x.MinVer != "" && minVer == "" {
		minVer = x.MinVer
	}

	// Union types: merge every member, carrying each member's microversion.
	for _, group := range [][]*base.SchemaProxy{s.OneOf, s.AnyOf, s.AllOf} {
		for _, m := range group {
			mv := minVer
			if r := refMinVer(m); r != "" {
				mv = r
			}
			flatten(m, prefix, mv, depth+1, cloneVisited(visited), out)
		}
	}

	// Arrays: descend into the item schema, keeping the same prefix.
	if s.Items != nil && s.Items.A != nil {
		flatten(s.Items.A, prefix, minVer, depth+1, visited, out)
	}

	if s.Properties == nil {
		return
	}
	required := map[string]bool{}
	for _, r := range s.Required {
		required[r] = true
	}
	for name, ps := range s.Properties.FromOldest() {
		child := ps.Schema()
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		// Only descend the outermost envelope (prefix == ""): the single
		// wrapping object ({"server": {...}}) or list wrapper ({"servers":
		// [...]}). Below that, a nested object or array-of-objects is a distinct
		// domain type that Gophercloud models as its own struct referenced by
		// this property, so we emit the property name as a leaf instead of
		// flattening its sub-fields into the parent's field pool.
		if prefix == "" {
			if child != nil && child.Properties != nil && child.Properties.Len() > 0 {
				// Wrapper object: descend, do not emit the wrapper itself.
				flatten(ps, path, minVer, depth+1, visited, out)
				continue
			}
			if child != nil && child.Items != nil && child.Items.A != nil {
				if is := child.Items.A.Schema(); is != nil && is.Properties != nil && is.Properties.Len() > 0 {
					flatten(ps, path, minVer, depth+1, visited, out)
					continue
				}
			}
		}
		f := model.Field{Name: name, Required: required[name], MinVer: minVer}
		if x := readX(childExt(child)); x.MinVer != "" {
			f.MinVer = x.MinVer
			f.MaxVer = x.MaxVer
		}
		if child != nil && child.Deprecated != nil {
			f.Deprecated = *child.Deprecated
		}
		*out = append(*out, f)
	}
}

// --- helpers ---

var refSuffix = regexp.MustCompile(`_(\d{2,})$`)

// mvFromRef derives a microversion from a component name like ".../Foo_294".
func mvFromRef(ref string) string {
	name := ref
	if i := strings.LastIndex(ref, "/"); i >= 0 {
		name = ref[i+1:]
	}
	m := refSuffix.FindStringSubmatch(name)
	if m == nil {
		return ""
	}
	digits := m[1]
	return digits[:1] + "." + digits[1:]
}

func refMinVer(m *base.SchemaProxy) string {
	if m == nil {
		return ""
	}
	return mvFromRef(m.GetReference())
}

func requestOneOf(rb *v3.RequestBody) []*base.SchemaProxy {
	if rb == nil || rb.Content == nil {
		return nil
	}
	mt := jsonMedia(rb.Content)
	if mt == nil || mt.Schema == nil {
		return nil
	}
	s := mt.Schema.Schema()
	if s == nil {
		return nil
	}
	return s.OneOf
}

func jsonMedia(content *orderedmap.Map[string, *v3.MediaType]) *v3.MediaType {
	if content == nil {
		return nil
	}
	if mt, ok := content.Get("application/json"); ok {
		return mt
	}
	// Fall back to the first available media type.
	for _, mt := range content.FromOldest() {
		return mt
	}
	return nil
}

func schemaExt(m *base.SchemaProxy) *orderedmap.Map[string, *yaml.Node] {
	if m == nil {
		return nil
	}
	return childExt(m.Schema())
}

func childExt(s *base.Schema) *orderedmap.Map[string, *yaml.Node] {
	if s == nil {
		return nil
	}
	return s.Extensions
}

func readX(ext *orderedmap.Map[string, *yaml.Node]) xOpenStack {
	var x xOpenStack
	if ext == nil {
		return x
	}
	node, ok := ext.Get("x-openstack")
	if !ok || node == nil {
		return x
	}
	_ = node.Decode(&x)
	return x
}

func verbs(item *v3.PathItem) map[string]*v3.Operation {
	return map[string]*v3.Operation{
		"GET":    item.Get,
		"PUT":    item.Put,
		"POST":   item.Post,
		"DELETE": item.Delete,
		"PATCH":  item.Patch,
		"HEAD":   item.Head,
	}
}

// stripPrefix removes the version prefix and normalises {param} to {}.
func stripPrefix(path, prefix string) string {
	if prefix != "" {
		path = strings.TrimPrefix(path, prefix)
	}
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return normalizePath(path)
}

// stripLeadingParam removes a single leading "{}" segment from an
// already-normalised path, e.g. "/{}/volumes" -> "/volumes" and "/{}" -> "/".
// Paths without a leading parameter segment are returned unchanged.
func stripLeadingParam(path string) string {
	if path == "/{}" {
		return "/"
	}
	if rest, ok := strings.CutPrefix(path, "/{}/"); ok {
		return "/" + rest
	}
	return path
}

// dedupeOps merges operations that share a Key(), unioning their fields and
// keeping the lowest microversion. It is used when stripProjectID collapses a
// project-scoped path onto its project-less twin.
func dedupeOps(ops []model.Operation) []model.Operation {
	byKey := map[string]int{}
	var out []model.Operation
	for _, op := range ops {
		if idx, ok := byKey[op.Key()]; ok {
			m := &out[idx]
			if lessVer(op.MinVer, m.MinVer) {
				m.MinVer = op.MinVer
			}
			m.RequestFields = dedupe(append(m.RequestFields, op.RequestFields...))
			m.ResponseFields = dedupe(append(m.ResponseFields, op.ResponseFields...))
			m.QueryParams = dedupe(append(m.QueryParams, op.QueryParams...))
			continue
		}
		byKey[op.Key()] = len(out)
		out = append(out, op)
	}
	return out
}

var pathParam = regexp.MustCompile(`\{[^}]+\}`)

func normalizePath(path string) string {
	path = pathParam.ReplaceAllString(path, "{}")
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}

func dedupe(fields []model.Field) []model.Field {
	seen := map[string]int{}
	var out []model.Field
	for _, f := range fields {
		if idx, ok := seen[f.Name]; ok {
			// Keep the lowest known microversion.
			if lessVer(f.MinVer, out[idx].MinVer) {
				out[idx].MinVer = f.MinVer
			}
			continue
		}
		seen[f.Name] = len(out)
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// lessVer reports whether microversion a is lower than b (empty == lowest).
func lessVer(a, b string) bool {
	if a == "" {
		return b != ""
	}
	if b == "" {
		return false
	}
	return a < b
}

func cloneVisited(v map[string]bool) map[string]bool {
	c := make(map[string]bool, len(v))
	for k := range v {
		c[k] = true
	}
	return c
}
