package hosts

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToHostListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
//
// Blazar accepts query parameters on this endpoint but ignores them. It returns every host.
type ListOpts struct{}

// ToHostListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of compute hosts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToHostListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return HostPage{pagination.SinglePageBase(r)}
	})
}
