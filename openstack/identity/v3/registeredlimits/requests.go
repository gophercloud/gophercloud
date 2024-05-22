package registeredlimits

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToRegisteredLimitListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Filters the response by a region ID.
	RegionID string `q:"region_id"`

	// Filters the response by a service ID.
	ServiceID string `q:"service_id"`

	// Filters the response by a resource name.
	ResourceName string `q:"resource_name"`
}

// ToRegisteredLimitListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRegisteredLimitListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the registered limits.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToRegisteredLimitListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RegisteredLimitPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// BatchCreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type BatchCreateOptsBuilder interface {
	ToRegisteredLimitsCreateMap() (map[string]any, error)
}

type CreateOpts struct {
	// RegionID is the ID of the region where the limit is applied.
	RegionID string `json:"region_id,omitempty"`

	// ServiceID is the ID of the service where the limit is applied.
	ServiceID string `json:"service_id" required:"true"`

	// Description of the limit.
	Description string `json:"description,omitempty"`

	// ResourceName is the name of the resource that the limit is applied to.
	ResourceName string `json:"resource_name" required:"true"`

	// DefaultLimit is the default limit.
	DefaultLimit int `json:"default_limit" required:"true"`
}

// BatchCreateOpts provides options used to create limits.
type BatchCreateOpts []CreateOpts

// ToRegisteredLimitsCreateMap formats a BatchCreateOpts into a create request.
func (opts BatchCreateOpts) ToRegisteredLimitsCreateMap() (map[string]any, error) {
	registered_limits := make([]map[string]any, len(opts))
	for i, registered_limit := range opts {
		registeredLimitMap, err := registered_limit.ToMap()
		if err != nil {
			return nil, err
		}
		registered_limits[i] = registeredLimitMap
	}
	return map[string]any{"registered_limits": registered_limits}, nil
}

func (opts CreateOpts) ToMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// BatchCreate creates new Limits.
func BatchCreate(ctx context.Context, client *gophercloud.ServiceClient, opts BatchCreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRegisteredLimitsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, rootURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves details on a single registered_limit, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, registeredLimitID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, registeredLimitID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToRegisteredLimitUpdateMap() (map[string]any, error)
}

// UpdateOpts represents parameters to update a domain.
type UpdateOpts struct {
	// Description of the registered_limit.
	Description *string `json:"description,omitempty"`

	// DefaultLimit is the override limit.
	DefaultLimit *int `json:"default_limit,omitempty"`

	// RegionID is the ID of the region where the limit is applied.
	RegionID string `json:"region_id,omitempty"`

	// ServiceID is the ID of the service where the limit is applied.
	ServiceID string `json:"service_id,omitempty"`

	// ResourceName is the name of the resource that the limit is applied to.
	ResourceName string `json:"resource_name,omitempty"`
	//Either service_id, resource_name, or region_id must be different than existing value otherwise it will raise 409.
}

// ToRegisteredLimitUpdateMap formats UpdateOpts into an update request.
func (opts UpdateOpts) ToRegisteredLimitUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "registered_limit")
}

// Update modifies the attributes of a registered limit.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRegisteredLimitUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, resourceURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a registered_limit.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, registeredLimitID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, registeredLimitID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
