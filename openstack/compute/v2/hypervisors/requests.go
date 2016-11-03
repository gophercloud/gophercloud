package hypervisors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToHypervisorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// ID of the server at which you want to set a marker.
	Marker string `q:"marker"`

	// Integer value for the limit of values to return.
	Limit int `q:"limit"`
}

// ToHypervisorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHypervisorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list hypervisors.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := hypervisorsListDetailURL(client)
	if opts != nil {
		query, err := opts.ToHypervisorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return HypervisorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
