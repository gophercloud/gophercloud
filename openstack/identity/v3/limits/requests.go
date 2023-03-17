package limits

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get retrieves details on a single limit, by ID.
func GetEnforcementModel(client *gophercloud.ServiceClient) (r EnforcementModelResult) {
	resp, err := client.Get(enforcementModelURL(client), &r.Body, nil)
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
	ToLimitsCreateMap() (map[string]interface{}, error)
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
func (opts BatchCreateOpts) ToLimitsCreateMap() (map[string]interface{}, error) {
	limits := make([]map[string]interface{}, len(opts))
	for i, limit := range opts {
		limitMap, err := limit.ToMap()
		if err != nil {
			return nil, err
		}
		limits[i] = limitMap
	}
	return map[string]interface{}{"limits": limits}, nil
}

func (opts CreateOpts) ToMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// BatchCreate creates new Limits.
func BatchCreate(client *gophercloud.ServiceClient, opts BatchCreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLimitsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(rootURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
