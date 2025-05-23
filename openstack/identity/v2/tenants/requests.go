package tenants

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToTenantListQuery() (string, error)
}

// ListOpts filters the Tenants that are returned by the List call.
type ListOpts struct {
	// Marker is the ID of the last Tenant on the previous page.
	Marker string `q:"marker"`

	// Limit specifies the page size.
	Limit int `q:"limit"`
}

// ToTenantListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTenantListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the Tenants to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToTenantListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
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

// CreateOptsBuilder enables extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTenantCreateMap() (map[string]any, error)
}

// ToTenantCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToTenantCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "tenant")
}

// Create is the operation responsible for creating new tenant.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTenantCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get requests details on a single tenant by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToTenantUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing
// tenant.
type UpdateOpts struct {
	// Name is the name of the tenant.
	Name string `json:"name,omitempty"`

	// Description is the description of the tenant.
	Description *string `json:"description,omitempty"`

	// Enabled sets the tenant status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ToTenantUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToTenantUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "tenant")
}

// Update is the operation responsible for updating exist tenants by their TenantID.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToTenantUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete is the operation responsible for permanently deleting a tenant.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
