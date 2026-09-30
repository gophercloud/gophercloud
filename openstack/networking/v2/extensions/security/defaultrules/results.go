package defaultrules

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// DefaultSecGroupRule represents a rule template used by Neutron to populate
// the rules of newly created security groups.
type DefaultSecGroupRule struct {
	// The UUID for this default security group rule.
	ID string

	// The direction in which the security group rule will be applied. The
	// only values allowed are "ingress" or "egress".
	Direction string

	// Description of the rule
	Description string `json:"description"`

	// Must be IPv4 or IPv6, and addresses represented in CIDR must match the
	// ingress or egress rules.
	EtherType string `json:"ethertype"`

	// The minimum port number in the range that will be matched by the
	// security group rule. If the protocol is TCP or UDP, this value must be
	// less than or equal to the value of the PortRangeMax attribute. If the
	// protocol is ICMP, this value must be an ICMP type.
	PortRangeMin int `json:"port_range_min"`

	// The maximum port number in the range that will be matched by the
	// security group rule. The PortRangeMin attribute constrains the
	// PortRangeMax attribute. If the protocol is ICMP, this value must be an
	// ICMP code.
	PortRangeMax int `json:"port_range_max"`

	// The protocol that will be matched by the security group rule. Valid
	// values are "tcp", "udp", "icmp" or an empty string.
	Protocol string

	// The remote address group ID to be associated with this security group
	// rule.
	RemoteAddressGroupID string `json:"remote_address_group_id"`

	// The remote group ID to be associated with this security group rule. The
	// special value "PARENT" references the security group the rule will
	// belong to once created from this template.
	RemoteGroupID string `json:"remote_group_id"`

	// The remote IP prefix to be associated with this security group rule.
	// This attribute matches the specified IP prefix as the source IP address
	// of the IP packet.
	RemoteIPPrefix string `json:"remote_ip_prefix"`

	// Whether this rule template is used in the default security group
	// created automatically for each new project.
	UsedInDefaultSG bool `json:"used_in_default_sg"`

	// Whether this rule template is used in security groups created by users,
	// other than the project default security group.
	UsedInNonDefaultSG bool `json:"used_in_non_default_sg"`

	// RevisionNumber optionally set via extensions/standard-attr-revisions
	RevisionNumber int `json:"revision_number"`

	// Timestamp when the rule was created
	CreatedAt time.Time `json:"-"`

	// Timestamp when the rule was last updated
	UpdatedAt time.Time `json:"-"`
}

func (r *DefaultSecGroupRule) UnmarshalJSON(b []byte) error {
	type tmp DefaultSecGroupRule

	// Support for older neutron time format
	var s1 struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339NoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339NoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s1)
	if err == nil {
		*r = DefaultSecGroupRule(s1.tmp)
		r.CreatedAt = time.Time(s1.CreatedAt)
		r.UpdatedAt = time.Time(s1.UpdatedAt)

		return nil
	}

	// Support for newer neutron time format
	var s2 struct {
		tmp
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = json.Unmarshal(b, &s2)
	if err != nil {
		return err
	}

	*r = DefaultSecGroupRule(s2.tmp)
	r.CreatedAt = time.Time(s2.CreatedAt)
	r.UpdatedAt = time.Time(s2.UpdatedAt)

	return nil
}

// DefaultSecGroupRulePage is the page returned by a pager when traversing over
// a collection of default security group rules.
type DefaultSecGroupRulePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of default security group
// rules has reached the end of a page and the pager seeks to traverse over a
// new one. In order to do this, it needs to construct the next page's URL.
func (r DefaultSecGroupRulePage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"default_security_group_rules_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a DefaultSecGroupRulePage struct is empty.
func (r DefaultSecGroupRulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractDefaultRules(r)
	return len(is) == 0, err
}

// ExtractDefaultRules accepts a Page struct, specifically a
// DefaultSecGroupRulePage struct, and extracts the elements into a slice of
// DefaultSecGroupRule structs. In other words, a generic collection is mapped
// into a relevant slice.
func ExtractDefaultRules(r pagination.Page) ([]DefaultSecGroupRule, error) {
	var s struct {
		DefaultSecGroupRules []DefaultSecGroupRule `json:"default_security_group_rules"`
	}
	err := (r.(DefaultSecGroupRulePage)).ExtractInto(&s)
	return s.DefaultSecGroupRules, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a default security
// group rule.
func (r commonResult) Extract() (*DefaultSecGroupRule, error) {
	var s struct {
		DefaultSecGroupRule *DefaultSecGroupRule `json:"default_security_group_rule"`
	}
	err := r.ExtractInto(&s)
	return s.DefaultSecGroupRule, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a DefaultSecGroupRule.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a DefaultSecGroupRule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
