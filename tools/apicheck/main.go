// Command apicheck compares OpenStack OpenAPI schemas against Gophercloud's
// implementation and reports coverage gaps.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gophercloud/gophercloud/tools/apicheck/config"
	"github.com/gophercloud/gophercloud/tools/apicheck/diff"
	"github.com/gophercloud/gophercloud/tools/apicheck/impl"
	"github.com/gophercloud/gophercloud/tools/apicheck/report"
	"github.com/gophercloud/gophercloud/tools/apicheck/spec"
)

func main() {
	cfgPath := flag.String("config", "apicheck.yaml", "path to apicheck config")
	specRoot := flag.String("spec-root", "", "override spec_root from config")
	gcRoot := flag.String("gc-root", "", "override gophercloud_root from config")
	only := flag.String("service", "", "limit to a single service name")
	format := flag.String("format", "markdown", "output format: markdown or json")
	dumpSpec := flag.Bool("dump-spec", false, "dump the parsed spec model as JSON and exit")
	dumpImpl := flag.Bool("dump-impl", false, "dump the parsed impl model as JSON and exit")
	baseline := flag.String("baseline", "", "path to a baseline JSON report; exit non-zero on operation regressions vs it")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fatal(err)
	}
	if *specRoot != "" {
		cfg.SpecRoot = *specRoot
	}
	if *gcRoot != "" {
		cfg.GophercloudRoot = *gcRoot
	}

	// gophercloud_root defaults to this repository's root, resolved relative to
	// the config file (tools/apicheck/apicheck.yaml -> ../..).
	if cfg.GophercloudRoot == "" {
		cfg.GophercloudRoot = filepath.Join(filepath.Dir(*cfgPath), "..", "..")
	}
	if cfg.SpecRoot == "" {
		fatal(fmt.Errorf("spec_root is not set: clone https://github.com/gtema/openstack " +
			"and pass its openstack_types/data directory via -spec-root (or the " +
			"APICHECK_SPEC_ROOT make variable)"))
	}

	var reports []diff.ServiceReport
	for _, svc := range cfg.Services {
		if *only != "" && svc.Name != *only {
			continue
		}

		specAPI, err := spec.LoadFromRoot(svc.Name, cfg.SpecRoot, svc.SpecFile, svc.PathPrefix, svc.StripProjectID)
		if err != nil {
			fatal(fmt.Errorf("%s: %w", svc.Name, err))
		}
		if *dumpSpec {
			dump(specAPI)
			continue
		}

		implRes, err := impl.Load(svc.Name, cfg.GophercloudRoot, svc.GophercloudPkg)
		if err != nil {
			fatal(fmt.Errorf("%s: %w", svc.Name, err))
		}
		if *dumpImpl {
			fmt.Fprintf(os.Stderr, "%s: %d impl operations, %d unresolved\n",
				svc.Name, len(implRes.API.Operations), len(implRes.Unresolved))
			dump(implRes.API)
			continue
		}

		reports = append(reports, diff.Compare(specAPI, implRes.API))
		if len(implRes.Unresolved) > 0 {
			fmt.Fprintf(os.Stderr, "%s: %d operations with unresolved URLs\n", svc.Name, len(implRes.Unresolved))
		}
	}

	if *dumpSpec || *dumpImpl {
		return
	}

	if *baseline != "" {
		checkBaseline(*baseline, reports)
		return
	}

	switch *format {
	case "json":
		if err := report.JSON(os.Stdout, reports); err != nil {
			fatal(err)
		}
	default:
		report.Markdown(os.Stdout, reports)
	}
}

// checkBaseline compares the current reports against a committed baseline and
// exits non-zero if any operation that was covered in the baseline is no longer
// covered. This is the intended CI regression gate.
func checkBaseline(path string, current []diff.ServiceReport) {
	data, err := os.ReadFile(path)
	if err != nil {
		fatal(fmt.Errorf("reading baseline %s: %w", path, err))
	}
	var base []diff.ServiceReport
	if err := json.Unmarshal(data, &base); err != nil {
		fatal(fmt.Errorf("parsing baseline %s: %w", path, err))
	}
	regs := diff.Regressions(base, current)
	if len(regs) == 0 {
		fmt.Println("apicheck: no operation coverage regressions vs baseline")
		return
	}
	fmt.Fprintf(os.Stderr, "apicheck: %d operation coverage regression(s) vs baseline:\n", len(regs))
	for _, r := range regs {
		fmt.Fprintf(os.Stderr, "  - [%s] %s (was %s)\n", r.Service, r.Key, r.Source)
	}
	os.Exit(1)
}

func dump(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "apicheck:", err)
	os.Exit(1)
}
