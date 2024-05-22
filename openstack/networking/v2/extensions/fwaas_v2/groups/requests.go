package groups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToFirewallGroupCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new firewall group.
type CreateOpts struct {
	ID                      string   `json:"id,omitempty"`
	TenantID                string   `json:"tenant_id,omitempty"`
	Name                    string   `json:"name,omitempty"`
	Description             string   `json:"description,omitempty"`
	IngressFirewallPolicyID string   `json:"ingress_firewall_policy_id,omitempty"`
	EgressFirewallPolicyID  string   `json:"egress_firewall_policy_id,omitempty"`
	AdminStateUp            *bool    `json:"admin_state_up,omitempty"`
	Ports                   []string `json:"ports,omitempty"`
	Shared                  *bool    `json:"shared,omitempty"`
	ProjectID               string   `json:"project_id,omitempty"`
}

// ToFirewallGroupCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToFirewallGroupCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "firewall_group")
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall group
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFirewallGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(ctx, rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToFirewallGroupUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating a firewall group.
type UpdateOpts struct {
	Name                    *string   `json:"name,omitempty"`
	Description             *string   `json:"description,omitempty"`
	IngressFirewallPolicyID *string   `json:"ingress_firewall_policy_id,omitempty"`
	EgressFirewallPolicyID  *string   `json:"egress_firewall_policy_id,omitempty"`
	AdminStateUp            *bool     `json:"admin_state_up,omitempty"`
	Ports                   *[]string `json:"ports,omitempty"`
	Shared                  *bool     `json:"shared,omitempty"`
}

// ToFirewallGroupUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToFirewallGroupUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "firewall_group")
}

// Update allows firewall groups to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFirewallGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Because of fwaas_v2 wait only UUID not string and base updateOpts has omitempty,
// only set nil allows firewall group policies to be unset.
// Two different functions, because can not specify both policy in one function.
// New functions needs new structs without omitempty.
// Separate function for BuildRequestBody is missing due to complication
// of code readability and bulkiness.

type RemoveIngressPolicyOpts struct {
	IngressFirewallPolicyID *string `json:"ingress_firewall_policy_id"`
}

func RemoveIngressPolicy(ctx context.Context, c *gophercloud.ServiceClient, id string) (r UpdateResult) {
	b, err := gophercloud.BuildRequestBody(RemoveIngressPolicyOpts{IngressFirewallPolicyID: nil}, "firewall_group")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type RemoveEgressPolicyOpts struct {
	EgressFirewallPolicyID *string `json:"egress_firewall_policy_id"`
}

func RemoveEgressPolicy(ctx context.Context, c *gophercloud.ServiceClient, id string) (r UpdateResult) {
	b, err := gophercloud.BuildRequestBody(RemoveEgressPolicyOpts{EgressFirewallPolicyID: nil}, "firewall_group")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular firewall group based on its unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(ctx, resourceURL(c, id), nil)
	return
}
