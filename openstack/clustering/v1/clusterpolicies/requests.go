package clusterpolicies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToClusterPolicyListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	Enabled *bool  `q:"enabled"`
	Name    string `q:"policy_name"`
	Type    string `q:"policy_type"`
	Sort    string `q:"sort"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of policies.
func List(client *gophercloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToClusterPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPolicyPage{pagination.SinglePageBase(r)}
	})
}
