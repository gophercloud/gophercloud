package addressgroups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToAddressGroupListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the address group attributes you want to see returned. SortKey allows
// you to sort by a particular network attribute. SortDir sets the direction,
// and is either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID          string   `q:"id"`
	Name        string   `q:"name"`
	Description string   `q:"description"`
	ProjectID   string   `q:"project_id"`
	Addresses   []string `q:"addresses"`
	Limit       int      `q:"limit"`
	Marker      string   `q:"marker"`
	SortKey     string   `q:"sort_key"`
	SortDir     string   `q:"sort_dir"`
}

// ToAddressGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAddressGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// address groups. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToAddressGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AddressGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAddressGroupCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new address group.
type CreateOpts struct {
	// The address group ID to associate with this address group.
	ID string `json:"id,omitempty"`

	// Human readable name for the address group (255 characters limit). Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human readable description for the address group (255 characters limit).
	Description string `json:"description,omitempty"`

	// Owner of the address group.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// Array of address. It supports both CIDR and IP range objects.
	// An example of addresses: [“132.168.4.12/24”, “132.168.5.12-132.168.5.24”, “2001::db8::f00/64”]
	Addresses []string `json:"addresses" required:"true"`
}

// ToAddressGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToAddressGroupCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "address_group")
}

// ToAddressesCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToAddressesCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create is an operation which creates a new address group and associates it
// with an existing address group (whose ID is specified in CreateOpts).
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAddressGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular address group based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular address group based on its
// unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update requests.
type UpdateOptsBuilder interface {
	ToAddressGroupUpdateMap() (map[string]any, error)
}

// UpdateOpts contains all the values needed to update an address group.
type UpdateOpts struct {
	// Human readable name for the address group (255 characters limit). Does not have to be unique.
	Name *string `json:"name,omitempty"`
	// Human readable description for the address group (255 characters limit).
	Description *string `json:"description,omitempty"`
}

// ToAddressGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToAddressGroupUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "address_group")
}

// Update will update a particular address group with a complete new set of data.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAddressGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateAddressesOpts will add or remove a list of particular addresses from
// an address group. It is used by both AddAddresses and RemoveAddresses
// requests.
type UpdateAddressesOpts struct {
	Addresses []string `json:"addresses" required:"true"`
}

// AddressesBuilder allows extensions to add additional parameters to the
// AddAddresses or RemoveAddresses requests.
type UpdateAddressesBuilder interface {
	ToUpdateAddressesMap() (map[string]any, error)
}

// ToAddressesCreateMap builds a request body from CreateOpts.
func (opts UpdateAddressesOpts) ToUpdateAddressesMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddAddresses will add IP addresses to a particular address group.
func AddAddresses(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateAddressesBuilder) (r AddAddressesResult) {
	b, err := opts.ToUpdateAddressesMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceAddAddressesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveAddresses will remove particular IP addresses from a particular address group.
func RemoveAddresses(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateAddressesBuilder) (r RemoveAddressesResult) {
	b, err := opts.ToUpdateAddressesMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceRemoveAddressesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
