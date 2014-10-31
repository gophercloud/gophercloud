package lb

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToLBListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	ChangesSince string `q:"changes-since"`
	Status       Status `q:"status"`
	NodeAddr     string `q:"nodeaddress"`
	Marker       string `q:"marker"`
	Limit        int    `q:"limit"`
}

// ToLBListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLBListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToLBListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return LBPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
