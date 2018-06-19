package receivers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts params
type ListOpts struct {
	Limit         int    `q:"limit"`
	Marker        string `q:"marker"`
	Sort          string `q:"sort"`
	GlobalProject string `q:"global_project"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	ClusterID     string `q:"cluster_id"`
	Action        string `q:"action"`
	User          string `q:"user"`
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToReceiverListQuery() (string, error)
}

// ToReceiverListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReceiverListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of cluster.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToReceiverListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ReceiverPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
