package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToServiceCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a service.
type CreateOpts struct {
	// Type is the type of the service.
	Type string `json:"type"`

	// Enabled is whether or not the service is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the service.
	Extra map[string]interface{} `json:"-"`
}

// ToServiceCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToServiceCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "service")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["service"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create adds a new service of the requested type to the catalog.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToServiceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// ListOptsBuilder enables extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToServiceListMap() (string, error)
}

// ListOpts provides options for filtering the List results.
type ListOpts struct {
	ServiceType string `q:"type"`
	PerPage     int    `q:"perPage"`
	Page        int    `q:"page"`
}

// ToServiceListMap builds a list query from the list options.
func (opts ListOpts) ToServiceListMap() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the services available to a specific user.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	u := listURL(client)
	if opts != nil {
		q, err := opts.ToServiceListMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		u += q
	}
	return pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, serviceID string) (r GetResult) {
	_, r.Err = client.Get(serviceURL(client, serviceID), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToServiceUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating a service.
type UpdateOpts struct {
	// Type is the type of the service.
	Type string `json:"type"`

	// Enabled is whether or not the service is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the service.
	Extra map[string]interface{} `json:"-"`
}

// ToServiceUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToServiceUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "service")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["service"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Service.
func Update(client *gophercloud.ServiceClient, serviceID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToServiceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, serviceID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete removes an existing service.
// It either deletes all associated endpoints, or fails until all endpoints
// are deleted.
func Delete(client *gophercloud.ServiceClient, serviceID string) (r DeleteResult) {
	_, r.Err = client.Delete(serviceURL(client, serviceID), nil)
	return
}
