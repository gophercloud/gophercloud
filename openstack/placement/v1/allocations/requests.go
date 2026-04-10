package allocations

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retrieves the allocations for a specific consumer by its UUID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, consumerUUID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, consumerUUID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ProviderAllocationsOpts specifies the resources to consume from a single resource
// provider in a write request.
type ProviderAllocationsOpts struct {
	// Resources maps resource class names to the integer amount to consume.
	Resources map[string]int `json:"resources"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToAllocationUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the allocation to be set for a consumer.
//
// This requires microversion 1.28 or later. Write operations on allocations
// using earlier microversions are not safe in concurrent environments.
//
// ConsumerGeneration must be set to nil when creating allocations for a new
// consumer (it serializes as JSON null, which signals to the server that no
// prior allocation is expected). For an existing consumer, set it to the
// generation value returned by a prior Get call. A mismatch causes a 409
// Conflict response, allowing the caller to retry safely.
type UpdateOpts struct {
	// Allocations maps resource provider UUIDs to the resources to consume.
	Allocations map[string]ProviderAllocationsOpts `json:"allocations"`

	// Required from microversion 1.8.
	ProjectID string `json:"project_id"`

	// Required from microversion 1.8.
	UserID string `json:"user_id"`

	// ConsumerGeneration must be nil for new consumers (serializes as null) or
	// the current generation for existing consumers.
	// See the UpdateOpts type documentation for details.
	ConsumerGeneration *int `json:"consumer_generation"`

	// Required from microversion 1.38.
	ConsumerType string `json:"consumer_type,omitempty"`
}

// ToAllocationUpdateMap constructs a request body from UpdateOpts.
func (opts UpdateOpts) ToAllocationUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update replaces all allocations for a consumer. The operation is atomic.
//
// Requires microversion 1.28 or later.
func Update(ctx context.Context, client *gophercloud.ServiceClient, consumerUUID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAllocationUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateURL(client, consumerUUID), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
