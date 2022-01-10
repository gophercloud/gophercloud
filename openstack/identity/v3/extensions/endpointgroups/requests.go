package endpointgroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get retrieves details on a single endpoint group, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToEndpointGroupListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Name filters the response by endpoint group name.
	Name string `q:"name"`
}

// ToEndpointGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToEndpointGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the endpoint groups
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToEndpointGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return EndpointGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
