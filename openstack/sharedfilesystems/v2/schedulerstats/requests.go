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

// ListOpts filters the storage pools returned by List.
type ListOpts struct {
	// Deprecated: ProjectID is part of the service endpoint, not a query filter.
	ProjectID string `json:"project_id,omitempty"`
	// The pool name for the back end.
	PoolName string `json:"pool_name" q:"pool"`
	// The host name for the back end.
	HostName string `json:"host_name" q:"host"`
	// The name of the back end.
	BackendName string `json:"backend_name" q:"backend"`
	// The capabilities for the storage back end.
	Capabilities string `json:"capabilities" q:"capabilities"`
	// The share type name or UUID. Allows filtering back end pools based on the extra-specs in the share type.
	// This requires microversion 2.23 or later.
	ShareType string `json:"share_type,omitempty" q:"share_type"`
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

// ListDetailOptsBuilder allows extensions to add additional parameters to the
// ListDetail request.
type ListDetailOptsBuilder interface {
	ToPoolsListQuery() (string, error)
}

// ListDetailOpts filters the storage pools returned by ListDetail.
type ListDetailOpts struct {
	// Deprecated: ProjectID is part of the service endpoint, not a query filter.
	ProjectID string `json:"project_id,omitempty"`
	// The pool name for the back end.
	PoolName string `json:"pool_name" q:"pool"`
	// The host name for the back end.
	HostName string `json:"host_name" q:"host"`
	// The name of the back end.
	BackendName string `json:"backend_name" q:"backend"`
	// The capabilities for the storage back end.
	Capabilities string `json:"capabilities" q:"capabilities"`
	// The share type name or UUID. Allows filtering back end pools based on the extra-specs in the share type.
	// This requires microversion 2.23 or later.
	ShareType string `json:"share_type,omitempty" q:"share_type"`
}

// ToPoolsListQuery formats a ListDetailOpts into a query string.
func (opts ListDetailOpts) ToPoolsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail makes a request against the API to list detailed pool information.
func ListDetail(client *gophercloud.ServiceClient, opts ListDetailOptsBuilder) pagination.Pager {
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
