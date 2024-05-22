package flavorprofiles

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorProfileListQuery() (string, error)
}

// ListOpts allows to manage the output of the request.
type ListOpts struct {
	// The fields that you want the server to return
	Fields []string `q:"fields"`
}

// ToFlavorProfileListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorProfileListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// FlavorProfiles. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToFlavorProfileListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorProfilePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToFlavorProfileCreateMap() (map[string]any, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name" required:"true"`

	// Providing the name of the provider supported by the Octavia installation.
	ProviderName string `json:"provider_name" required:"true"`

	// Providing the json string containing the flavor metadata.
	FlavorData string `json:"flavor_data" required:"true"`
}

// ToFlavorProfileCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFlavorProfileCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "flavorprofile")
}

// Create is and operation which add a new FlavorProfile into the database.
// CreateResult will be returned.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFlavorProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular FlavorProfile based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToFlavorProfileUpdateMap() (map[string]any, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Providing the name of the provider supported by the Octavia installation.
	ProviderName string `json:"provider_name,omitempty"`

	// Providing the json string containing the flavor metadata.
	FlavorData string `json:"flavor_data,omitempty"`
}

// ToFlavorProfileUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToFlavorProfileUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "flavorprofile")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update is an operation which modifies the attributes of the specified
// FlavorProfile.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToFlavorProfileUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular FlavorProfile based on its
// unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
