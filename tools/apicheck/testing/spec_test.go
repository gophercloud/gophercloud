package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/tools/apicheck/model"
	"github.com/gophercloud/gophercloud/tools/apicheck/spec"
)

func loadSpec(t *testing.T) *model.API {
	t.Helper()
	api, err := spec.Load("fake", "fixtures/spec/widget.yaml", "", false)
	if err != nil {
		t.Fatalf("spec.Load: %v", err)
	}
	return api
}

// TestSpecEnvelopeFlatten verifies that the flattener descends the single
// request/response envelope but treats nested domain objects and
// array-of-objects as leaf fields rather than flattening their sub-fields into
// the parent's pool.
func TestSpecEnvelopeFlatten(t *testing.T) {
	api := loadSpec(t)

	op, ok := opByKey(api, "POST /widgets")
	if !ok {
		t.Fatalf("POST /widgets not found; ops: %+v", api.Operations)
	}

	// Envelope leaves and the microversion union member are present.
	for _, want := range []string{"name", "description", "tags", "metadata", "size"} {
		if _, ok := fieldByName(op.RequestFields, want); !ok {
			t.Errorf("expected request field %q: got %+v", want, op.RequestFields)
		}
	}
	// The "widget" envelope itself is not emitted.
	if _, ok := fieldByName(op.RequestFields, "widget"); ok {
		t.Errorf("envelope wrapper 'widget' must not be emitted as a field")
	}
	// Sub-fields of nested objects/arrays must NOT be flattened into the pool.
	for _, unwanted := range []string{"key", "value", "foo"} {
		if _, ok := fieldByName(op.RequestFields, unwanted); ok {
			t.Errorf("nested sub-field %q must not be flattened into the parent pool", unwanted)
		}
	}

	// The name field is required.
	if name, _ := fieldByName(op.RequestFields, "name"); !name.Required {
		t.Errorf("expected name to be required")
	}
	// size comes from WidgetCreate_21 -> microversion 2.1.
	if size, _ := fieldByName(op.RequestFields, "size"); size.MinVer != "2.1" {
		t.Errorf("expected size min-ver 2.1, got %q", size.MinVer)
	}
}

// TestSpecQueryAndResponse verifies query params (with microversion) and the
// list-envelope response flattening.
func TestSpecQueryAndResponse(t *testing.T) {
	api := loadSpec(t)

	op, ok := opByKey(api, "GET /widgets")
	if !ok {
		t.Fatalf("GET /widgets not found; ops: %+v", api.Operations)
	}
	if _, ok := fieldByName(op.QueryParams, "name"); !ok {
		t.Errorf("expected query param name")
	}
	limit, ok := fieldByName(op.QueryParams, "limit")
	if !ok {
		t.Fatalf("expected query param limit")
	}
	if limit.MinVer != "1.5" {
		t.Errorf("expected limit min-ver 1.5, got %q", limit.MinVer)
	}
	// List envelope {"widgets": [...]} collapses to item leaves.
	for _, want := range []string{"id", "name"} {
		if _, ok := fieldByName(op.ResponseFields, want); !ok {
			t.Errorf("expected response field %q: got %+v", want, op.ResponseFields)
		}
	}
}
