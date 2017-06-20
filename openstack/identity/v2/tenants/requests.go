package tenants

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts filters the Tenants that are returned by the List call.
type ListOpts struct {
	// Marker is the ID of the last Tenant on the previous page.
	Marker string `q:"marker"`
	// Limit specifies the page size.
	Limit int `q:"limit"`
}

// List enumerates the Tenants to which the current token has access.
func List(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		q, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q.String()
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TenantPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOpts represents the options needed when creating new tenant.
type CreateOpts struct {
	// Name is the name of the tenant.
	Name string `json:"name" required:"true"`
	// Description is the description of the tenant.
	Description string `json:"description,omitempty"`
	// Enabled sets the tenant status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
type CreateOptsBuilder interface {
	ToTenantCreateMap() (map[string]interface{}, error)
}

// ToTenantCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToTenantCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "tenant")
}

// Create is the operation responsible for creating new tenant.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTenantCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}
