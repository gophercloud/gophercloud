package services

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToServiceListQuery() (string, error)
}

// ListOpts holds options for listing Services.
type ListOpts struct {
	// The pool name for the back end.
	ProjectID string `json:"project_id,omitempty"`
	// The service host name.
	Host string `json:"host"`
	// The service binary name. Default is the base name of the executable.
	Binary string `json:"binary"`
	// The availability zone.
	Zone string `json:"zone"`
	// The current state of the service. A valid value is up or down.
	State string `json:"state"`
	// The service status, which is enabled or disabled.
	Status string `json:"status"`
}

// ToServiceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServiceListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list services.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToServiceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.SinglePageBase(r)}
	})
}
