package defaultrules

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToDefaultSecGroupRuleListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the default security group rule attributes you want to see returned. SortKey
// allows you to sort by a particular network attribute. SortDir sets the
// direction, and is either `asc' or `desc'. Marker and Limit are used for
// pagination.
type ListOpts struct {
	ID                   string `q:"id"`
	Description          string `q:"description"`
	Direction            string `q:"direction"`
	EtherType            string `q:"ethertype"`
	PortRangeMax         int    `q:"port_range_max"`
	PortRangeMin         int    `q:"port_range_min"`
	Protocol             string `q:"protocol"`
	RemoteAddressGroupID string `q:"remote_address_group_id"`
	RemoteGroupID        string `q:"remote_group_id"`
	RemoteIPPrefix       string `q:"remote_ip_prefix"`
	UsedInDefaultSG      *bool  `q:"used_in_default_sg"`
	UsedInNonDefaultSG   *bool  `q:"used_in_non_default_sg"`
	Limit                int    `q:"limit"`
	Marker               string `q:"marker"`
	SortKey              string `q:"sort_key"`
	SortDir              string `q:"sort_dir"`
}

// ToDefaultSecGroupRuleListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDefaultSecGroupRuleListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// default security group rules. It accepts a ListOpts struct, which allows you
// to filter and sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToDefaultSecGroupRuleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return DefaultSecGroupRulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToDefaultSecGroupRuleCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new default security
// group rule. Rules created through this API are templates used by Neutron to
// populate the rules of newly created security groups; they are not applied to
// any existing security group.
type CreateOpts struct {
	// Must be either "ingress" or "egress": the direction in which the
	// security group rule will be applied.
	Direction rules.RuleDirection `json:"direction" required:"true"`

	// String description of each rule, optional.
	Description string `json:"description,omitempty"`

	// Must be "IPv4" or "IPv6". Defaults to "IPv4" when not set.
	EtherType rules.RuleEtherType `json:"ethertype,omitempty"`

	// The maximum port number in the range that will be matched by the
	// security group rule. The PortRangeMin attribute constrains the
	// PortRangeMax attribute. If the protocol is ICMP, this value must be an
	// ICMP code.
	PortRangeMax *int `json:"port_range_max,omitempty"`

	// The minimum port number in the range that will be matched by the
	// security group rule. If the protocol is TCP or UDP, this value must be
	// less than or equal to the value of the PortRangeMax attribute. If the
	// protocol is ICMP, this value must be an ICMP type.
	PortRangeMin *int `json:"port_range_min,omitempty"`

	// The protocol that will be matched by the security group rule. Valid
	// values are "tcp", "udp", "icmp" or an empty string.
	Protocol rules.RuleProtocol `json:"protocol,omitempty"`

	// The remote address group ID to be associated with this security group
	// rule. You can specify either RemoteAddressGroupID, RemoteGroupID, or
	// RemoteIPPrefix.
	RemoteAddressGroupID string `json:"remote_address_group_id,omitempty"`

	// The remote group ID to be associated with this security group rule. You
	// can specify either RemoteAddressGroupID, RemoteGroupID or
	// RemoteIPPrefix. The special value "PARENT" can be used to reference the
	// security group the rule will belong to once created from this template.
	RemoteGroupID string `json:"remote_group_id,omitempty"`

	// The remote IP prefix to be associated with this security group rule.
	// You can specify either RemoteAddressGroupID, RemoteGroupID or
	// RemoteIPPrefix. This attribute matches the specified IP prefix as the
	// source IP address of the IP packet.
	RemoteIPPrefix string `json:"remote_ip_prefix,omitempty"`

	// Whether this rule template should be used in the default security group
	// created automatically for each new project. Defaults to false.
	UsedInDefaultSG *bool `json:"used_in_default_sg,omitempty"`

	// Whether this rule template should be used in security groups created by
	// users, other than the project default security group. Defaults to true.
	UsedInNonDefaultSG *bool `json:"used_in_non_default_sg,omitempty"`
}

// ToDefaultSecGroupRuleCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToDefaultSecGroupRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "default_security_group_rule")
}

// Create is an operation which adds a new default security group rule
// template that will be used by Neutron when creating the rules of new
// security groups.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToDefaultSecGroupRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular default security group rule based on its unique
// ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular default security group rule
// based on its unique ID. Existing security groups are not affected.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
