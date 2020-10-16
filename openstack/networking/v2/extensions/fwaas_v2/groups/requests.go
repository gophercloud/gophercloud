package groups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToGroupListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the firewall group attributes you want to see returned. SortKey allows you
// to sort by a particular firewall group attribute. SortDir sets the direction,
// and is either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	TenantID                string    `q:"tenant_id"`
	Name                    string    `q:"name"`
	Description             string    `q:"description"`
	IngressFirewallPolicyID string    `q:"ingress_firewall_policy_id"`
	EgressFirewallPolicyID  string    `q:"egress_firewall_policy_id"`
	AdminStateUp            *bool     `q:"admin_state_up"`
	Ports                   *[]string `q:"ports"`
	Status                  string    `q:"status"`
	ID                      string    `q:"id"`
	Shared                  *bool     `q:"shared"`
	ProjectID               string    `q:"project_id"`
	Limit                   int       `q:"limit"`
	Marker                  string    `q:"marker"`
	SortKey                 string    `q:"sort_key"`
	SortDir                 string    `q:"sort_dir"`
}

// ToGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// firewall groups. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
//
// Default group settings return only those firewall groups that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)

	if opts != nil {
		query, err := opts.ToGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a particular firewall group based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}
