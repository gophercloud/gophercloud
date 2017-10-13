package domains

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToDomainListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Enabled filters the response by enabled domains.
	Enabled *bool `q:"enabled"`

	// Name filters the response by domain name.
	Name string `q:"name"`
}

// ToDomainListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDomainListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the domains to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDomainListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DomainPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single domain, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
