package firewall_groups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	//"fmt"
)

// FirewallGroup is an OpenStack firewall_group.
type FirewallGroup struct {
	ID           	string `json:"id"`
	Name         	string `json:"name"`
	Description  	string `json:"description"`
	AdminStateUp 	bool   `json:"admin_state_up"`
	Status       	string `json:"status"`
	IngressPolicyID	string `json:"ingress_firewall_policy_id"`
	EgressPolicyID  string `json:"egress_firewall_policy_id"`
	TenantID     	string `json:"tenant_id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a firewall.
func (r commonResult) Extract() (*FirewallGroup, error) {
	var s FirewallGroup
	//fmt.Printf("Extracting %s.\n", r.PrettyPrintJSON())
	err := r.ExtractInto(&s)
	//fmt.Printf("Extracted %+v.\n", s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "firewall_group")
}

func ExtractFirewallGroupsInto(r pagination.Page, v interface{}) error {
	return r.(FirewallGroupPage).Result.ExtractIntoSlicePtr(v, "firewall_groups")
}

// FirewallPage is the page returned by a pager when traversing over a
// collection of firewalls.
type FirewallGroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of firewalls has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FirewallGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"firewalls_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FirewallPage struct is empty.
func (r FirewallGroupPage) IsEmpty() (bool, error) {
	is, err := ExtractFirewallGroups(r)
	return len(is) == 0, err
}

// ExtractFirewalls accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFirewallGroups(r pagination.Page) ([]FirewallGroup, error) {
	var s []FirewallGroup
	err := ExtractFirewallGroupsInto(r, &s)
	return s, err
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}
