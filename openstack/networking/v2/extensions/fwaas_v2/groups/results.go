package groups

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Group is a firewall group.
type Group struct {
	ID                      string   `json:"id"`
	TenantID                string   `json:"tenant_id"`
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	IngressFirewallPolicyID string   `json:"ingress_firewall_policy_id"`
	EgressFirewallPolicyID  string   `json:"egress_firewall_policy_id"`
	AdminStateUp            bool     `json:"admin_state_up"`
	Ports                   []string `json:"ports"`
	Status                  string   `json:"status"`
	Shared                  bool     `json:"shared"`
	ProjectID               string   `json:"project_id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a firewall group.
func (r commonResult) Extract() (*Group, error) {
	var s struct {
		Group *Group `json:"firewall_group"`
	}
	err := r.ExtractInto(&s)
	return s.Group, err
}

// GroupPage is the page returned by a pager when traversing over a
// collection of firewall groups.
type GroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of firewall groups has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r GroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"firewall_groups_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a GroupPage struct is empty.
func (r GroupPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractGroups(r)
	return len(is) == 0, err
}

// ExtractGroups accepts a Page struct, specifically a GroupPage struct,
// and extracts the elements into a slice of Group structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractGroups(r pagination.Page) ([]Group, error) {
	var s struct {
		Groups []Group `json:"firewall_groups"`
	}
	err := (r.(GroupPage)).ExtractInto(&s)
	return s.Groups, err
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
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
