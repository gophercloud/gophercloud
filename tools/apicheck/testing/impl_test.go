package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/tools/apicheck/impl"
	"github.com/gophercloud/gophercloud/tools/apicheck/model"
)

// loadFixture runs the AST extractor over the fakeservice fixture module.
func loadFixture(t *testing.T) *model.API {
	t.Helper()
	res, err := impl.Load("fake", "fixtures", "fakeservice")
	if err != nil {
		t.Fatalf("impl.Load: %v", err)
	}
	if len(res.Unresolved) != 0 {
		t.Fatalf("unexpected unresolved URLs: %+v", res.Unresolved)
	}
	return res.API
}

func opByKey(api *model.API, key string) (model.Operation, bool) {
	for _, op := range api.Operations {
		if op.Key() == key {
			return op, true
		}
	}
	return model.Operation{}, false
}

func fieldByName(fields []model.Field, name string) (model.Field, bool) {
	for _, f := range fields {
		if f.Name == name {
			return f, true
		}
	}
	return model.Field{}, false
}

// TestCreateStructBody covers a BuildRequestBody struct body, including the
// required tag and a json:"-" field surfaced as ManualHandled.
func TestCreateStructBody(t *testing.T) {
	api := loadFixture(t)

	op, ok := opByKey(api, "POST /widgets")
	if !ok {
		t.Fatalf("POST /widgets not found; ops: %+v", api.Operations)
	}
	if op.Source != "fakeservice.Create" {
		t.Errorf("expected source fakeservice.Create, got %q", op.Source)
	}

	name, ok := fieldByName(op.RequestFields, "name")
	if !ok {
		t.Fatalf("request field name not found: %+v", op.RequestFields)
	}
	if !name.Required {
		t.Errorf("expected name to be required")
	}
	for _, want := range []string{"description", "size"} {
		if _, ok := fieldByName(op.RequestFields, want); !ok {
			t.Errorf("request field %q not found", want)
		}
	}
	// json:"-" field is carried by its Go name and marked ManualHandled.
	secret, ok := fieldByName(op.RequestFields, "Secret")
	if !ok {
		t.Fatalf("expected json:\"-\" field Secret to be present as ManualHandled")
	}
	if !secret.ManualHandled {
		t.Errorf("expected Secret to be ManualHandled")
	}
}

// TestListQueryParams covers NewPager list detection and q: query opts,
// including a q:"-" field that must be excluded.
func TestListQueryParams(t *testing.T) {
	api := loadFixture(t)

	op, ok := opByKey(api, "GET /widgets")
	if !ok {
		t.Fatalf("GET /widgets not found; ops: %+v", api.Operations)
	}
	for _, want := range []string{"name", "limit"} {
		if _, ok := fieldByName(op.QueryParams, want); !ok {
			t.Errorf("query param %q not found: %+v", want, op.QueryParams)
		}
	}
	if _, ok := fieldByName(op.QueryParams, "Internal"); ok {
		t.Errorf("q:\"-\" field Internal must be excluded from query params")
	}
}

// TestNestedURLResolution covers a url helper that delegates to another helper
// (deleteURL -> resourceURL -> ServiceURL).
func TestNestedURLResolution(t *testing.T) {
	api := loadFixture(t)
	if _, ok := opByKey(api, "DELETE /widgets/{}"); !ok {
		t.Fatalf("DELETE /widgets/{} not found; ops: %+v", api.Operations)
	}
}

// TestManualMapActionBody covers a manual map[string]any action body whose
// single key names the action.
func TestManualMapActionBody(t *testing.T) {
	api := loadFixture(t)
	op, ok := opByKey(api, "POST /widgets/{}/action:reboot")
	if !ok {
		t.Fatalf("action reboot not found; ops: %+v", api.Operations)
	}
	if op.Action != "reboot" {
		t.Errorf("expected action reboot, got %q", op.Action)
	}
}
