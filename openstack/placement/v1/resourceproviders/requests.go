package resourceproviders

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToResourceProviderListQuery() (string, error)
}

// ListOpts allows the filtering resource providers. Filtering is achieved by
// passing in struct field values that map to the resource provider attributes
// you want to see returned.
type ListOpts struct {
	// Name is the name of the resource provider to filter the list
	Name string `q:"name"`

	// UUID is the uuid of the resource provider to filter the list
	UUID string `q:"uuid"`

	// MemberOf is a string representing aggregate uuids to filter or exclude from the list
	MemberOf string `q:"member_of"`

	// Resources is a comma-separated list of string indicating an amount of resource
	// of a specified class that a provider must have the capacity and availability to serve
	Resources string `q:"resources"`

	// InTree is a string that represents a resource provider UUID.  The returned resource
	// providers will be in the same provider tree as the specified provider.
	InTree string `q:"in_tree"`

	// Required is comma-delimited list of string trait names.
	Required string `q:"required"`
}

// ToResourceProviderListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToResourceProviderListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list resource providers.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := resourceProvidersListURL(client)

	if opts != nil {
		query, err := opts.ToResourceProviderListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourceProvidersPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToResourceProviderCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a resource provider.
type CreateOpts struct {
	Name string `json:"name"`
	UUID string `json:"uuid,omitempty"`
	// The UUID of the immediate parent of the resource provider.
	// Available in version >= 1.14
	ParentProviderUUID string `json:"parent_provider_uuid,omitempty"`
}

// ToResourceProviderCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToResourceProviderCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create makes a request against the API to create a resource provider
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToResourceProviderCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, resourceProvidersListURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the resource provider associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, resourceProviderID string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, resourceProviderID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a specific resource provider based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, resourceProviderID string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, resourceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToResourceProviderUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a resource provider.
type UpdateOpts struct {
	Name *string `json:"name,omitempty"`
	// Available in version >= 1.37. It can be set to any existing provider UUID
	// except to providers that would cause a loop. Also it can be set to null
	// to transform the provider to a new root provider. This operation needs to
	// be used carefully. Moving providers can mean that the original rules used
	// to create the existing resource allocations may be invalidated by that move.
	ParentProviderUUID *string `json:"parent_provider_uuid,omitempty"`
}

// ToResourceProviderUpdateMap constructs a request body from UpdateOpts.
func (opts UpdateOpts) ToResourceProviderUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update makes a request against the API to create a resource provider
func Update(ctx context.Context, client *gophercloud.ServiceClient, resourceProviderID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToResourceProviderUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURL(client, resourceProviderID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetUsages(ctx context.Context, client *gophercloud.ServiceClient, resourceProviderID string) (r GetUsagesResult) {
	resp, err := client.Get(ctx, getResourceProviderUsagesURL(client, resourceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetInventories(ctx context.Context, client *gophercloud.ServiceClient, resourceProviderID string) (r GetInventoriesResult) {
	resp, err := client.Get(ctx, getResourceProviderInventoriesURL(client, resourceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetAllocations(ctx context.Context, client *gophercloud.ServiceClient, resourceProviderID string) (r GetAllocationsResult) {
	resp, err := client.Get(ctx, getResourceProviderAllocationsURL(client, resourceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetTraits(ctx context.Context, client *gophercloud.ServiceClient, resourceProviderID string) (r GetTraitsResult) {
	resp, err := client.Get(ctx, getResourceProviderTraitsURL(client, resourceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
