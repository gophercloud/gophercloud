package schedulerstats

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToPoolsListQuery() (string, error)
}

// ListOpts controls the view of data returned (e.g globally or per project).
type ListOpts struct {
	// The pool name for the back end.
	PoolName string `q:"pool"`
	// The host name for the back end.
	HostName string `q:"host"`
	// The name of the back end.
	BackendName string `q:"backend"`
	// The capabilities for the storage back end.
	Capabilities string `q:"capabilities"`
	// The share type name or UUID. Allows filtering back end pools based on the extra-specs in the share type.
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
