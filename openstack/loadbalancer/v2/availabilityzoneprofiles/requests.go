package availabilityzoneprofiles

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToAvailabilityZoneProfileListQuery() (string, error)
}

// ListOpts allows to manage the output of the request.
type ListOpts struct {
	// The name of the availability zone profile to filter by.
	Name string `q:"name"`
	// The provider name of the availability zone profile to filter by.
	ProviderName string `q:"provider_name"`
	// The fields that you want the server to return
	Fields []string `q:"fields"`
}

// ToAvailabilityZoneProfileListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAvailabilityZoneProfileListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// AvailabilityZoneProfiles. It accepts a ListOpts struct, which allows you to
// filter and sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToAvailabilityZoneProfileListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AvailabilityZoneProfilePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAvailabilityZoneProfileCreateMap() (map[string]any, error)
}

// CreateOpts is the common options struct used in this package's create
// operation.
type CreateOpts struct {
	// Human-readable name for the avaialability zone profile.
	// Does not have to be unique.
	Name string `json:"name" required:"true"`

	// Providing the name of the provider supported by the Octavia installation.
	ProviderName string `json:"provider_name" required:"true"`

	// Providing the json string containing the availability zone metadata.
	AvailabilityZoneData string `json:"availability_zone_data" required:"true"`
}

// ToAvailabilityZoneProfileCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToAvailabilityZoneProfileCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "availability_zone_profile")
}

// Create is and operation which add a new AvailabilityZoneProfile into the database.
// CreateResult will be returned.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAvailabilityZoneProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular AvailabilityZoneProfile based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// update request.
type UpdateOptsBuilder interface {
	ToAvailabiltyZoneProfileUpdateMap() (map[string]any, error)
}

// UpdateOpts is the common options struct used in this package's update
// operation.
type UpdateOpts struct {
	// Human-readable name for the availability zone profile.
	// Does not have to be unique.
	Name *string `json:"name,omitempty"`

	// Providing the name of the provider supported by the Octavia installation.
	ProviderName *string `json:"provider_name,omitempty"`

	// Providing the json string containing the availability zone metadata.
	AvailabiltyZoneData *string `json:"availability_zone_data,omitempty"`
}

// ToAvailabiltyZoneProfileUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToAvailabiltyZoneProfileUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "availability_zone_profile")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update is an operation which modifies the attributes of the specified
// AvailabilityZoneProfile.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAvailabiltyZoneProfileUpdateMap()
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

// Delete will permanently delete a particular AvailabiltyZoneProfile based on
// its unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
