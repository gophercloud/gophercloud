package limits

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Get retrieves details on a single limit, by ID.
func GetEnforcementModel(ctx context.Context, client *gophercloud.ServiceClient) (r EnforcementModelResult) {
	resp, err := client.Get(ctx, enforcementModelURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToLimitListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Filters the response by a region ID.
	RegionID string `q:"region_id"`

	// Filters the response by a project ID.
	ProjectID string `q:"project_id"`

	// Filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Filters the response by a service ID.
	ServiceID string `q:"service_id"`

	// Filters the response by a resource name.
	ResourceName string `q:"resource_name"`
}

// ToLimitListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLimitListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the limits.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToLimitListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return LimitPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// BatchCreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type BatchCreateOptsBuilder interface {
	ToLimitsCreateMap() (map[string]any, error)
}

type CreateOpts struct {
	// RegionID is the ID of the region where the limit is applied.
	RegionID string `json:"region_id,omitempty"`

	// ProjectID is the ID of the project where the limit is applied.
	ProjectID string `json:"project_id,omitempty"`

	// DomainID is the ID of the domain where the limit is applied.
	DomainID string `json:"domain_id,omitempty"`

	// ServiceID is the ID of the service where the limit is applied.
	ServiceID string `json:"service_id" required:"true"`

	// Description of the limit.
	Description string `json:"description,omitempty"`

	// ResourceName is the name of the resource that the limit is applied to.
	ResourceName string `json:"resource_name" required:"true"`

	// ResourceLimit is the override limit.
	ResourceLimit int `json:"resource_limit"`
}

// BatchCreateOpts provides options used to create limits.
type BatchCreateOpts []CreateOpts

// ToLimitsCreateMap formats a BatchCreateOpts into a create request.
func (opts BatchCreateOpts) ToLimitsCreateMap() (map[string]any, error) {
	limits := make([]map[string]any, len(opts))
	for i, limit := range opts {
		limitMap, err := limit.ToMap()
		if err != nil {
			return nil, err
		}
		limits[i] = limitMap
	}
	return map[string]any{"limits": limits}, nil
}

func (opts CreateOpts) ToMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// BatchCreate creates new Limits.
func BatchCreate(ctx context.Context, client *gophercloud.ServiceClient, opts BatchCreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLimitsCreateMap()
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

// Get retrieves details on a single limit, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, limitID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, limitID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToLimitUpdateMap() (map[string]any, error)
}

// UpdateOpts represents parameters to update a domain.
type UpdateOpts struct {
	// Description of the limit.
	Description *string `json:"description,omitempty"`

	// ResourceLimit is the override limit.
	ResourceLimit *int `json:"resource_limit,omitempty"`
}

// ToLimitUpdateMap formats UpdateOpts into an update request.
func (opts UpdateOpts) ToLimitUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "limit")
}

// Update modifies the attributes of a limit.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToLimitUpdateMap()
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

// Delete deletes a limit.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, limitID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, limitID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
