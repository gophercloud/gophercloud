package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/tools/apicheck/diff"
	"github.com/gophercloud/gophercloud/tools/apicheck/model"
)

// TestDiffCoverage checks operation matching, field set-difference, the
// ManualHandled exemption, and unmatched-impl reporting.
func TestDiffCoverage(t *testing.T) {
	spec := &model.API{Service: "fake", Operations: []model.Operation{
		{Service: "fake", Method: "POST", Path: "/widgets", RequestFields: []model.Field{
			{Name: "name", Required: true},
			{Name: "description"},
			{Name: "token"},
			{Name: "size", MinVer: "2.1"},
		}},
		{Service: "fake", Method: "GET", Path: "/widgets/{}"}, // missing op
	}}
	impl := &model.API{Service: "fake", Operations: []model.Operation{
		{Service: "fake", Method: "POST", Path: "/widgets", Source: "fakeservice.Create", RequestFields: []model.Field{
			{Name: "name", Required: true},
			{Name: "description"},
			{Name: "token", ManualHandled: true}, // hand-serialised; must not be a gap
		}},
		{Service: "fake", Method: "DELETE", Path: "/widgets/{}", Source: "fakeservice.Delete"}, // extra impl
	}}

	rep := diff.Compare(spec, impl)

	if rep.SpecOps != 2 || rep.ImplementedOp != 1 {
		t.Fatalf("expected 2 spec ops / 1 implemented, got %d / %d", rep.SpecOps, rep.ImplementedOp)
	}
	if got := rep.OperationCoverage(); got != 50 {
		t.Errorf("expected 50%% coverage, got %.0f", got)
	}

	var post *diff.OperationGap
	for i := range rep.Gaps {
		if rep.Gaps[i].Key == "POST /widgets" {
			post = &rep.Gaps[i]
		}
	}
	if post == nil {
		t.Fatal("POST /widgets gap not found")
	}
	// Only "size" is genuinely missing: "secret"/"Secret" is ManualHandled.
	if len(post.MissingReq) != 1 || post.MissingReq[0].Name != "size" {
		t.Fatalf("expected only 'size' missing, got %+v", post.MissingReq)
	}
	if post.MissingReq[0].MinVer != "2.1" {
		t.Errorf("expected size min-ver 2.1, got %q", post.MissingReq[0].MinVer)
	}

	if len(rep.ExtraImpl) != 1 {
		t.Errorf("expected 1 unmatched impl op, got %+v", rep.ExtraImpl)
	}
}

// TestRegressions checks that only operations that lost coverage are flagged,
// while newly-unimplemented spec operations are not.
func TestRegressions(t *testing.T) {
	baseline := []diff.ServiceReport{{
		Service: "fake",
		Gaps: []diff.OperationGap{
			{Key: "POST /widgets", Implemented: true, Source: "fakeservice.Create"},
			{Key: "GET /widgets", Implemented: true, Source: "fakeservice.List"},
			{Key: "GET /widgets/{}", Implemented: false}, // was already a gap
		},
	}}
	current := []diff.ServiceReport{{
		Service: "fake",
		Gaps: []diff.OperationGap{
			{Key: "POST /widgets", Implemented: true, Source: "fakeservice.Create"},
			{Key: "GET /widgets", Implemented: false}, // regressed: coverage lost
			{Key: "GET /widgets/{}", Implemented: false},
			{Key: "DELETE /widgets/{}", Implemented: false}, // new spec op, not a regression
		},
	}}

	regs := diff.Regressions(baseline, current)
	if len(regs) != 1 {
		t.Fatalf("expected exactly 1 regression, got %+v", regs)
	}
	if regs[0].Key != "GET /widgets" || regs[0].Service != "fake" {
		t.Errorf("unexpected regression: %+v", regs[0])
	}
	if regs[0].Source != "fakeservice.List" {
		t.Errorf("expected baseline source carried through, got %q", regs[0].Source)
	}
}
