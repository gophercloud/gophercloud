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
	Enabled bool   `q:"enabled"`
	Name    string `q:"policy_name,omitempty"`
	Type    string `q:"policy_type,omitempty"`
	Sort    string `q:"sort,omitempty"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details of a single policy. Use ExtractPolicy to convert its
// result into a Node.
func Get(client *gophercloud.ServiceClient, clusterID string, policyID string) (r GetResult) {
	_, r.Err = client.Get(getDetailURL(client, clusterID, policyID), &r.Body, nil)
	return
}

// ListDetail instructs OpenStack to provide a list of policies.
func ListDetail(client *gophercloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToClusterPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
