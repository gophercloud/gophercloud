package resourceproviders

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type ResourceProviderLinks struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// ResourceProvider are entities which provider consumable inventory of one or more classes of resource
type ResourceProvider struct {
	// Generation is a consistent view marker that assists with the management of concurrent resource provider updates.
	Generation int `json:"generation"`

	// UUID of a resource provider.
	UUID string `json:"uuid"`

	// Links is a list of links associated with one resource provider.
	Links []ResourceProviderLinks `json:"links"`

	// Name of one resource provider.
	Name string `json:"name"`

	// The ParentProviderUUID contains the UUID of the immediate parent of the resource provider.
	// Requires microversion 1.14 or above
	ParentProviderUUID string `json:"parent_provider_uuid"`

	// The RootProviderUUID contains the read-only UUID of the top-most provider in this provider tree.
	// Requires microversion 1.14 or above
	RootProviderUUID string `json:"root_provider_uuid"`
}

type ResourceProviderUsage struct {
	ResourceProviderGeneration int            `json:"resource_provider_generation"`
	Usages                     map[string]int `json:"usages"`
}

type Inventory struct {
	AllocationRatio float32 `json:"allocation_ratio"`
	MaxUnit         int     `json:"max_unit"`
	MinUnit         int     `json:"min_unit"`
	Reserved        int     `json:"reserved"`
	StepSize        int     `json:"step_size"`
	Total           int     `json:"total"`
}

type Allocation struct {
	Resources map[string]int `json:"resources"`
}

type ResourceProviderInventories struct {
	ResourceProviderGeneration int                  `json:"resource_provider_generation"`
	Inventories                map[string]Inventory `json:"inventories"`
}

type ResourceProviderAllocations struct {
	ResourceProviderGeneration int                   `json:"resource_provider_generation"`
	Allocations                map[string]Allocation `json:"allocations"`
}

type ResourceProviderTraits struct {
	ResourceProviderGeneration int      `json:"resource_provider_generation"`
	Traits                     []string `json:"traits"`
}

// resourceProviderResult is the response of a base ResourceProvider result.
type resourceProviderResult struct {
	gophercloud.Result
}

// Extract interpets any resourceProviderResult-base result as a ResourceProvider.
func (r resourceProviderResult) Extract() (*ResourceProvider, error) {
	var s ResourceProvider
	err := r.ExtractInto(&s)

	return &s, err
}

// CreateResult is the result of a Create operation. Call its Extract
// method to interpret it as a ResourceProvider.
type CreateResult struct {
	resourceProviderResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult represents the result of a create operation. Call its Extract
// method to interpret it as a ResourceProvider.
type GetResult struct {
	resourceProviderResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a ResourceProvider.
type UpdateResult struct {
	resourceProviderResult
}

// ResourceProvidersPage contains a single page of all resource providers from a List call.
type ResourceProvidersPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines if a ResourceProvidersPage contains any results.
func (page ResourceProvidersPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	resourceProviders, err := ExtractResourceProviders(page)
	return len(resourceProviders) == 0, err
}

// ExtractResourceProviders returns a slice of ResourceProvider from a List operation.
func ExtractResourceProviders(r pagination.Page) ([]ResourceProvider, error) {
	var s struct {
		ResourceProviders []ResourceProvider `json:"resource_providers"`
	}
	err := (r.(ResourceProvidersPage)).ExtractInto(&s)
	return s.ResourceProviders, err
}

// GetUsagesResult is the response of a Get usage operations. Call its Extract method
// to interpret it as a ResourceProviderUsage.
type GetUsagesResult struct {
	gophercloud.Result
}

// Extract interprets a GetUsagesResult as a ResourceProviderUsage.
func (r GetUsagesResult) Extract() (*ResourceProviderUsage, error) {
	var s ResourceProviderUsage
	err := r.ExtractInto(&s)
	return &s, err
}

// GetInventoriesResult is the response of a Get inventories operations. Call its Extract method
// to interpret it as a ResourceProviderInventories.
type GetInventoriesResult struct {
	gophercloud.Result
}

// Extract interprets a GetInventoriesResult as a ResourceProviderInventories.
func (r GetInventoriesResult) Extract() (*ResourceProviderInventories, error) {
	var s ResourceProviderInventories
	err := r.ExtractInto(&s)
	return &s, err
}

// GetAllocationsResult is the response of a Get allocations operations. Call its Extract method
// to interpret it as a ResourceProviderAllocations.
type GetAllocationsResult struct {
	gophercloud.Result
}

// Extract interprets a GetAllocationsResult as a ResourceProviderAllocations.
func (r GetAllocationsResult) Extract() (*ResourceProviderAllocations, error) {
	var s ResourceProviderAllocations
	err := r.ExtractInto(&s)
	return &s, err
}

// GetTraitsResult is the response of a Get traits operations. Call its Extract method
// to interpret it as a ResourceProviderTraits.
type GetTraitsResult struct {
	gophercloud.Result
}

// Extract interprets a GetTraitsResult as a ResourceProviderTraits.
func (r GetTraitsResult) Extract() (*ResourceProviderTraits, error) {
	var s ResourceProviderTraits
	err := r.ExtractInto(&s)
	return &s, err
}
