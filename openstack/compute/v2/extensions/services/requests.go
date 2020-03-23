package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request.
type ListOptsBuilder interface {
	ToServicesListQuery() (string, error)
}

// ListOpts represents options to list services.
type ListOpts struct {
	Binary string `q:"binary"`
	Host   string `q:"host"`
}

// ToServicesListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServicesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list services.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToServicesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.SinglePageBase(r)}
	})
}
