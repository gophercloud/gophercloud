package allocations

import "github.com/gophercloud/gophercloud/v2"

// ProviderAllocations represents the per-provider portion of an allocations response.
type ProviderAllocations struct {
	Generation int `json:"generation"`

	// Resources maps resource class names to the integer amount consumed from the provider.
	Resources map[string]int `json:"resources"`
}

// Allocations represents the allocations for a single consumer.
// The Allocations field maps resource provider UUIDs to ProviderAllocations,
// describing how much of each resource class is consumed from each provider.
type Allocations struct {
	// Allocations maps resource provider UUIDs to the resources consumed from each.
	Allocations map[string]ProviderAllocations `json:"allocations"`

	// Available from microversion 1.12.
	// Will be absent when listing allocations for a consumer UUID that has no allocations.
	ProjectID *string `json:"project_id"`

	// Available from microversion 1.12.
	// Will be absent when listing allocations for a consumer UUID that has no allocations.
	UserID *string `json:"user_id"`

	// Available from microversion 1.28.
	ConsumerGeneration *int `json:"consumer_generation"`

	// Available from microversion 1.38.
	ConsumerType *string `json:"consumer_type"`
}

// GetResult is the result of a Get operation. Call its Extract method
// to interpret it as an Allocations.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as Allocations.
func (r GetResult) Extract() (*Allocations, error) {
	var s Allocations
	err := r.ExtractInto(&s)
	return &s, err
}

// UpdateResult is the result of an Update operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type UpdateResult struct {
	gophercloud.ErrResult
}

// DeleteResult is the result of a Delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ManageResult is the result of a Manage operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type ManageResult struct {
	gophercloud.ErrResult
}
