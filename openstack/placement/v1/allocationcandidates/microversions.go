package allocationcandidates

import "github.com/gophercloud/gophercloud/v2/pagination"

// AllocationCandidates110 represents the response from a List allocation
// candidates request for microversions 1.10-1.11.
// In these versions the allocations field is an array rather than a dictionary.
type AllocationCandidates110 struct {
	// AllocationRequests is a list of objects that contain information
	// for creating a later allocation claim request.
	AllocationRequests []AllocationRequest110 `json:"allocation_requests"`

	// ProviderSummaries is a dictionary keyed by resource provider UUID of
	// inventory/capacity information for providers in the allocation_requests.
	ProviderSummaries map[string]ProviderSummary110 `json:"provider_summaries"`
}

// AllocationRequest110 represents a single allocation request for
// microversions 1.10-1.11 where allocations is an array.
type AllocationRequest110 struct {
	// Allocations is a list of allocation resources per provider.
	Allocations []AllocationRequest110Resource `json:"allocations"`
}

// AllocationRequest110Resource represents a single provider allocation
// within a 1.10-1.11 allocation request.
type AllocationRequest110Resource struct {
	// ResourceProvider contains the UUID of the resource provider.
	ResourceProvider AllocationRequest110ResourceProvider `json:"resource_provider"`

	// Resources is a dictionary of resource class names to the amount requested.
	Resources map[string]int `json:"resources"`
}

// AllocationRequest110ResourceProvider contains the UUID of a resource provider.
type AllocationRequest110ResourceProvider struct {
	UUID string `json:"uuid"`
}

// ProviderSummary110 represents a provider summary for microversions 1.10-1.11
// which only includes resources (no traits or parent/root UUIDs).
type ProviderSummary110 struct {
	// Resources is a dictionary of resource class names to capacity/usage info.
	Resources map[string]ProviderSummaryResource `json:"resources"`
}

// ExtractAllocationCandidates110 interprets an AllocationCandidatesPage as AllocationCandidates110 (microversions 1.10-1.11).
func ExtractAllocationCandidates110(r pagination.Page) (*AllocationCandidates110, error) {
	var s AllocationCandidates110
	err := (r.(AllocationCandidatesPage)).ExtractInto(&s)
	return &s, err
}
