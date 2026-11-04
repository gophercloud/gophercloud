// Package diff compares a spec-side API against an impl-side API and reports
// coverage gaps.
package diff

import (
	"sort"

	"github.com/gophercloud/gophercloud/tools/apicheck/model"
)

// FieldGap is a spec field absent from the implementation.
type FieldGap struct {
	Name   string `json:"name"`
	MinVer string `json:"min_ver,omitempty"`
	MaxVer string `json:"max_ver,omitempty"`
}

// OperationGap describes the coverage of a single spec operation.
type OperationGap struct {
	Key          string     `json:"key"`
	Method       string     `json:"method"`
	Path         string     `json:"path"`
	Action       string     `json:"action,omitempty"`
	MinVer       string     `json:"min_ver,omitempty"`
	Implemented  bool       `json:"implemented"`
	Source       string     `json:"source,omitempty"`
	MissingReq   []FieldGap `json:"missing_request_fields,omitempty"`
	MissingQuery []FieldGap `json:"missing_query_params,omitempty"`
	MissingResp  []FieldGap `json:"missing_response_fields,omitempty"`
}

// ServiceReport is the diff outcome for one service.
type ServiceReport struct {
	Service       string         `json:"service"`
	SpecOps       int            `json:"spec_operations"`
	ImplementedOp int            `json:"implemented_operations"`
	Gaps          []OperationGap `json:"operations"`
	// ExtraImpl lists impl operations with no matching spec operation (possible
	// path-mapping mismatch or a Gophercloud extension not in the schema).
	ExtraImpl []string `json:"extra_impl,omitempty"`
}

// OperationCoverage returns implemented / total as a percentage.
func (r ServiceReport) OperationCoverage() float64 {
	if r.SpecOps == 0 {
		return 100
	}
	return 100 * float64(r.ImplementedOp) / float64(r.SpecOps)
}

// Compare diffs spec against impl.
func Compare(spec, impl *model.API) ServiceReport {
	implByKey := map[string]model.Operation{}
	for _, op := range impl.Operations {
		implByKey[op.Key()] = op
	}
	specKeys := map[string]bool{}

	rep := ServiceReport{Service: spec.Service, SpecOps: len(spec.Operations)}
	for _, sop := range spec.Operations {
		specKeys[sop.Key()] = true
		gap := OperationGap{
			Key:    sop.Key(),
			Method: sop.Method,
			Path:   sop.Path,
			Action: sop.Action,
			MinVer: sop.MinVer,
		}
		iop, ok := implByKey[sop.Key()]
		if !ok {
			rep.Gaps = append(rep.Gaps, gap)
			continue
		}
		gap.Implemented = true
		gap.Source = iop.Source
		rep.ImplementedOp++
		gap.MissingReq = missing(sop.RequestFields, iop.RequestFields)
		gap.MissingQuery = missing(sop.QueryParams, iop.QueryParams)
		gap.MissingResp = missing(sop.ResponseFields, iop.ResponseFields)
		rep.Gaps = append(rep.Gaps, gap)
	}

	for _, iop := range impl.Operations {
		if !specKeys[iop.Key()] {
			rep.ExtraImpl = append(rep.ExtraImpl, iop.Key()+"  ("+iop.Source+")")
		}
	}
	sort.Strings(rep.ExtraImpl)
	return rep
}

// Regression is an operation that was covered in a baseline report but is no
// longer covered in the current run.
type Regression struct {
	Service string `json:"service"`
	Key     string `json:"key"`
	Source  string `json:"baseline_source,omitempty"`
}

// Regressions returns operations that were implemented in the baseline but are
// no longer implemented in current. Only operation-level regressions are
// reported: field-level differences are dominated by spec updates adding fields
// and are too noisy to gate on. New spec operations that are simply unimplemented
// are not regressions.
func Regressions(baseline, current []ServiceReport) []Regression {
	curCovered := map[string]map[string]bool{} // service -> key -> implemented
	for _, r := range current {
		m := map[string]bool{}
		for _, g := range r.Gaps {
			if g.Implemented {
				m[g.Key] = true
			}
		}
		curCovered[r.Service] = m
	}

	var out []Regression
	for _, r := range baseline {
		cur := curCovered[r.Service]
		for _, g := range r.Gaps {
			if !g.Implemented {
				continue
			}
			if cur == nil || !cur[g.Key] {
				out = append(out, Regression{Service: r.Service, Key: g.Key, Source: g.Source})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Service != out[j].Service {
			return out[i].Service < out[j].Service
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// missing returns spec fields whose name is not present on the impl side.
// Fields marked ManualHandled on the impl side count as present (they exist but
// are serialised by hand and cannot be verified by tag).
func missing(spec, impl []model.Field) []FieldGap {
	have := map[string]bool{}
	for _, f := range impl {
		have[f.Name] = true
	}
	var gaps []FieldGap
	for _, f := range spec {
		if have[f.Name] {
			continue
		}
		gaps = append(gaps, FieldGap{Name: f.Name, MinVer: f.MinVer, MaxVer: f.MaxVer})
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i].Name < gaps[j].Name })
	return gaps
}
