// Package config describes how Gophercloud packages map onto the OpenAPI
// schema files, and holds per-service overrides for the coverage check.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level apicheck configuration.
type Config struct {
	// SpecRoot is the path to the gtema/openstack schema data directory
	// (…/openstack_types/data). Left unset in the committed config; supply it via
	// -spec-root (or the APICHECK_SPEC_ROOT make variable).
	SpecRoot string `yaml:"spec_root"`
	// GophercloudRoot is the path to the Gophercloud module root. Left unset in
	// the committed config, in which case it defaults to the repository root;
	// override with -gc-root.
	GophercloudRoot string `yaml:"gophercloud_root"`
	// Services enumerates the services to check.
	Services []Service `yaml:"services"`
}

// Service maps one Gophercloud service/version onto one OpenAPI schema file.
type Service struct {
	// Name is the logical service key used in reports (e.g. "compute").
	Name string `yaml:"name"`
	// GophercloudPkg is the package path relative to the Gophercloud module
	// root (e.g. "openstack/compute/v2").
	GophercloudPkg string `yaml:"gophercloud_pkg"`
	// SpecFile is the OpenAPI file relative to SpecRoot
	// (e.g. "compute/v2.yaml").
	SpecFile string `yaml:"spec_file"`
	// PathPrefix is the leading path segment(s) present in the OpenAPI paths
	// but not in Gophercloud's ServiceURL segments (e.g. "/v2.1"). Stripped
	// before matching.
	PathPrefix string `yaml:"path_prefix"`
	// StripProjectID collapses a leading project-scope path parameter after the
	// version prefix. Some services (notably block-storage/Cinder) list every
	// path twice — once project-scoped ("/v3/{project_id}/volumes") and once
	// project-less ("/v3/volumes") — because the project ID is legacy in the URL
	// and modern deployments carry it in the service-catalog endpoint (which is
	// what Gophercloud relies on). With this set, "/{project_id}/volumes"
	// normalises to "/volumes" so the two variants dedupe to a single operation.
	StripProjectID bool `yaml:"strip_project_id"`
	// IgnoreOperations lists operation keys (e.g. "GET /os-hosts") that are
	// intentionally not implemented and should be excluded from gap counts.
	IgnoreOperations []string `yaml:"ignore_operations"`
}

// Load reads a Config from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &c, nil
}
