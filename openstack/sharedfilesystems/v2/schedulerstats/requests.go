package schedulerstats

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List and ListDetail requests.
type ListOptsBuilder interface {
	ToPoolsListQuery() (string, error)
}

// ListOpts filters the storage pools returned by List and ListDetail.
type ListOpts struct {
	// Pool is a regular expression to match against the pool name.
	Pool string `q:"pool"`
	// Host is a regular expression to match against the host name.
	Host string `q:"host"`
	// Backend is a regular expression to match against the back-end name.
	Backend string `q:"backend"`
	// ShareType is the share type name or UUID. It filters pools by matching
	// their capabilities against the share type's extra specs.
	// This requires microversion 2.23 or later.
	ShareType string `q:"share_type"`
}

// ToPoolsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPoolsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list pool information.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := poolsListURL(client)
	if opts != nil {
		query, err := opts.ToPoolsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.SinglePageBase(r)}
	})
}

// ListDetail makes a request against the API to list detailed pool information.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := poolsListDetailURL(client)
	if opts != nil {
		query, err := opts.ToPoolsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.SinglePageBase(r)}
	})
}
