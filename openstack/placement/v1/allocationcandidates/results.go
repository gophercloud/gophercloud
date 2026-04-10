package allocationcandidates

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/v2/pagination"
)

// AllocationCandidates represents the response from a List allocation
// candidates request for microversions 1.12 and above.
type AllocationCandidates struct {
	// AllocationRequests is a list of objects that contain information
	// for creating a later allocation claim request.
	AllocationRequests []AllocationRequest `json:"allocation_requests"`

	// ProviderSummaries is a dictionary keyed by resource provider UUID of
	// inventory/capacity information for providers in the allocation_requests.
	ProviderSummaries map[string]ProviderSummary `json:"provider_summaries"`
}

// AllocationRequest represents a single allocation request within the
// allocation candidates response.
type AllocationRequest struct {
	// Allocations is a dictionary of resource allocations keyed by
	// resource provider UUID.
	Allocations map[string]AllocationRequestResource `json:"allocations"`

	// Mappings is a dictionary associating request group suffixes with a
	// list of UUIDs identifying the resource providers that satisfied each group.
	// Available in version >= 1.34.
	Mappings *map[string][]string `json:"mappings"`
}

// AllocationRequestResource represents the resources requested from a
// single resource provider within an allocation request.
type AllocationRequestResource struct {
	// Resources is a dictionary of resource class names to the amount requested.
	Resources map[string]int `json:"resources"`
}

// ProviderSummary represents the summary of a resource provider's inventory
// and traits.
type ProviderSummary struct {
	// Resources is a dictionary of resource class names to capacity/usage info.
	Resources map[string]ProviderSummaryResource `json:"resources"`

	// Traits is a list of traits the resource provider has associated with it.
	// Available in version >= 1.17.
	Traits *[]string `json:"traits"`

	// ParentProviderUUID is the UUID of the immediate parent of the resource provider.
	// Available in version >= 1.29.
	ParentProviderUUID *string `json:"parent_provider_uuid"`

	// RootProviderUUID is the UUID of the top-most provider in this provider tree.
	// Available in version >= 1.29.
	RootProviderUUID *string `json:"root_provider_uuid"`
}

// ProviderSummaryResource represents the capacity and usage of a single
// resource class for a provider.
type ProviderSummaryResource struct {
	// Capacity is the amount of the resource that the provider can accommodate.
	Capacity int `json:"capacity"`

	// Used is the amount of the resource that has been already allocated.
	Used int `json:"used"`
}

// AllocationCandidatesPage is the page returned from a List call.
type AllocationCandidatesPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines if an AllocationCandidatesPage contains any results.
// It avoids full deserialization so that it works across all microversions,
func (page AllocationCandidatesPage) IsEmpty() (bool, error) {
	var s struct {
		AllocationRequests []json.RawMessage `json:"allocation_requests"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return false, err
	}
	return len(s.AllocationRequests) == 0, nil
}

// ExtractAllocationCandidates interprets an AllocationCandidatesPage as AllocationCandidates (microversion 1.12+).
func ExtractAllocationCandidates(r pagination.Page) (*AllocationCandidates, error) {
	var s AllocationCandidates
	err := (r.(AllocationCandidatesPage)).ExtractInto(&s)
	return &s, err
}
