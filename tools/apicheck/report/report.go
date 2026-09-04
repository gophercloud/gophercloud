// Package report renders diff results as Markdown or JSON.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/gophercloud/gophercloud/tools/apicheck/diff"
)

// Markdown writes a human-readable coverage report.
func Markdown(w io.Writer, reports []diff.ServiceReport) {
	fmt.Fprintln(w, "# Gophercloud API coverage report")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Service | Operations | Covered | Coverage |")
	fmt.Fprintln(w, "|---|---|---|---|")
	for _, r := range reports {
		fmt.Fprintf(w, "| %s | %d | %d | %.0f%% |\n",
			r.Service, r.SpecOps, r.ImplementedOp, r.OperationCoverage())
	}
	fmt.Fprintln(w)

	for _, r := range reports {
		fmt.Fprintf(w, "## %s (%.0f%% of %d operations)\n\n", r.Service, r.OperationCoverage(), r.SpecOps)

		var missingOps, partial []diff.OperationGap
		for _, g := range r.Gaps {
			switch {
			case !g.Implemented:
				missingOps = append(missingOps, g)
			case len(g.MissingReq)+len(g.MissingQuery)+len(g.MissingResp) > 0:
				partial = append(partial, g)
			}
		}

		if len(missingOps) > 0 {
			fmt.Fprintf(w, "### Missing operations (%d)\n\n", len(missingOps))
			for _, g := range missingOps {
				fmt.Fprintf(w, "- `%s`%s\n", label(g), minVer(g.MinVer))
			}
			fmt.Fprintln(w)
		}

		if len(partial) > 0 {
			fmt.Fprintf(w, "### Partially-covered operations (%d)\n\n", len(partial))
			for _, g := range partial {
				fmt.Fprintf(w, "- `%s` (%s)\n", label(g), g.Source)
				printFields(w, "request field", g.MissingReq)
				printFields(w, "query param", g.MissingQuery)
				printFields(w, "response field", g.MissingResp)
			}
			fmt.Fprintln(w)
		}

		if len(r.ExtraImpl) > 0 {
			fmt.Fprintf(w, "<details><summary>Unmatched impl operations (%d) — likely path-mapping mismatches</summary>\n\n", len(r.ExtraImpl))
			for _, e := range r.ExtraImpl {
				fmt.Fprintf(w, "- `%s`\n", e)
			}
			fmt.Fprintln(w, "\n</details>")
		}
		fmt.Fprintln(w)
	}
}

func printFields(w io.Writer, kind string, gaps []diff.FieldGap) {
	if len(gaps) == 0 {
		return
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i].Name < gaps[j].Name })
	parts := make([]string, 0, len(gaps))
	for _, g := range gaps {
		s := g.Name
		if g.MinVer != "" {
			s += " (" + g.MinVer + ")"
		}
		parts = append(parts, s)
	}
	fmt.Fprintf(w, "  - missing %ss: ", kind)
	for i, p := range parts {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprint(w, p)
	}
	fmt.Fprintln(w)
}

func label(g diff.OperationGap) string {
	s := g.Method + " " + g.Path
	if g.Action != "" {
		s += " :" + g.Action
	}
	return s
}

func minVer(v string) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf(" _(min-ver %s)_", v)
}

// JSON writes the raw diff results.
func JSON(w io.Writer, reports []diff.ServiceReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}
