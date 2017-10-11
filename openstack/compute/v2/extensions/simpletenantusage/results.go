package simpletenantusage

import (
	"code.comcast.com/onecloud/gophercloud"
	"code.comcast.com/onecloud/gophercloud/pagination"
)

// TenantUsage is a set of usage information about a tenant over the sampling window
type TenantUsage struct {
	// ServerUsages is an array of ServerUsage maps
	ServerUsages []ServerUsage `json:"server_usages"`

	// Start is the beginning time to calculate usage statistics on compute and storage resources
	Start string `json:"start"`

	// Stop is the ending time to calculate usage statistics on compute and storage resources
	Stop string `json:"stop"`

	// TenantID is the ID of the tenant whose usage is being reported on
	TenantID string `json:"tenant_id"`

	// TotalHours is the total duration that servers exist (in hours)
	TotalHours float64 `json:"total_hours"`

	// TotalLocalGBUsage multiplies the server disk size (in GiB) by hours the server exists, and then adding that all together for each server
	TotalLocalGBUsage float64 `json:"total_local_gb_usage"`

	// TotalMemoryMBUsage multiplies the server memory size (in MB) by hours the server exists, and then adding that all together for each server
	TotalMemoryMBUsage float64 `json:"total_memory_mb_usage"`

	// TotalVCPUsUsage multiplies the number of virtual CPUs of the server by hours the server exists, and then adding that all together for each server
	TotalVCPUsUsage float64 `json:"total_vcpus_usage"`
}

// ServerUsage is a detailed set of information about a specific instance inside a tenant
type ServerUsage struct {
	// EndedAt is the date and time when the server was deleted
	EndedAt string `json:"ended_at"`

	// Flavor is the display name of a flavor
	Flavor string `json:"flavor"`

	// Hours is the duration that the server exists in hours
	Hours float64 `json:"hours"`

	// InstanceID is the UUID of the instance
	InstanceID string `json:"instance_id"`

	// LocalGB is the sum of the root disk size of the server and the ephemeral disk size of it (in GiB)
	LocalGB int `json:"local_gb"`

	// MemoryMB is the memory size of the server (in MB)
	MemoryMB int `json:"memory_mb"`

	// Name is the name assigned to the server when it was created
	Name string `json:"name"`

	// StartedAt is the date and time when the server was started
	StartedAt string `json:"started_at"`

	// State is the VM power state
	State string `json:"state"`

	// TenantID is the UUID of the tenant in a multi-tenancy cloud
	TenantID string `json:"tenant_id"`

	// Uptime is the uptime of the server in seconds
	Uptime int `json:"uptime"`

	// VCPUs is the number of virtual CPUs that the server uses
	VCPUs int `json:"vcpus"`
}

// TenantUsageLinks is an optional item with links pertaining to usage
type TenantUsageLinks struct {
	// HREF is the URI link to the usage documentation
	HREF string

	// Rel is the link relation for this link, indicating the relation or purpose of the link
	Rel string
}

// SimpleTenantUsagePage stores a single, only page of SimpleTenantUsage results from a List call.
type SimpleTenantUsagePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a SimpleTenantUsagePage is empty.
func (page SimpleTenantUsagePage) IsEmpty() (bool, error) {
	ks, err := ExtractSimpleTenantUsages(page)
	return ks == nil, err
}

// Type to specifically indicate Simple Tenant Usage results.
type simpleTenantUsageResult struct {
	gophercloud.Result
}

// ExtractSimpleTenantUsage is a function that attempts to interpret any SimpleTenantUsage resource response as a SimpleTenantUsage struct.
// The difference between ExtractSimpleTenantUsage and ExtractSimpleTenantUsages is that when a tenant ID is provided the JSON is
// "tenant_usage" (singular) which is a struct, otherwise it is "tenant_usages" (plural) which is an array of structs.
func ExtractSimpleTenantUsage(page pagination.Page) (*[]TenantUsage, error) {
	var s struct {
		TenantUsages     []TenantUsage      `json:"tenant_usage"`
		TenantUsageLinks []TenantUsageLinks `json:"tenant_usage_links"`
	}
	err := (page.(SimpleTenantUsagePage)).ExtractInto(&s)
	return &s.TenantUsages, err
}

// ExtractSimpleTenantUsages is a function that attempts to interpret any SimpleTenantUsage resource response as a SimpleTenantUsage struct.
// The difference between ExtractSimpleTenantUsage and ExtractSimpleTenantUsages is that when a tenant ID is provided the JSON is
// "tenant_usage" (singular) which is a struct, otherwise it is "tenant_usages" (plural) which is an array of structs.
func ExtractSimpleTenantUsages(page pagination.Page) (*[]TenantUsage, error) {
	var s struct {
		TenantUsages     []TenantUsage      `json:"tenant_usages"`
		TenantUsageLinks []TenantUsageLinks `json:"tenant_usage_links"`
	}
	err := (page.(SimpleTenantUsagePage)).ExtractInto(&s)
	return &s.TenantUsages, err
}

// GetResult is the response from a Get operation. Call its Extract function to interpret it
// as a SimpleTenantUsage.
type GetResult struct {
	simpleTenantUsageResult
}
