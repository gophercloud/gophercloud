package rules

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Rule represents a firewall rule
type Rule struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name,omitempty"`
	Description          string    `json:"description,omitempty"`
	Protocol             string    `json:"protocol"`
	Action               string    `json:"action"`
	IPVersion            int       `json:"ip_version,omitempty"`
	SourceIPAddress      string    `json:"source_ip_address,omitempty"`
	DestinationIPAddress string    `json:"destination_ip_address,omitempty"`
	SourcePort           string    `json:"source_port,omitempty"`
	DestinationPort      string    `json:"destination_port,omitempty"`
	Shared               bool      `json:"shared,omitempty"`
	Enabled              bool      `json:"enabled,omitempty"`
	FirewallPolicyID     []string  `json:"firewall_policy_id"`
	TenantID             string    `json:"tenant_id"`
	ProjectID            string    `json:"project_id"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"`
}

func (r *Rule) UnmarshalJSON(b []byte) error {
	type tmp Rule

	// Support for older neutron time format
	var s1 struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339NoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339NoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s1)
	if err == nil {
		*r = Rule(s1.tmp)
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

	*r = Rule(s2.tmp)
	r.CreatedAt = time.Time(s2.CreatedAt)
	r.UpdatedAt = time.Time(s2.UpdatedAt)

	return nil
}

// RulePage is the page returned by a pager when traversing over a
// collection of firewall rules.
type RulePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of firewall rules has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r RulePage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"firewall_rules_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RulePage struct is empty.
func (r RulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractRules(r)
	return len(is) == 0, err
}

// ExtractRules accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRules(r pagination.Page) ([]Rule, error) {
	var s struct {
		Rules []Rule `json:"firewall_rules"`
	}
	err := (r.(RulePage)).ExtractInto(&s)
	return s.Rules, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a firewall rule.
func (r commonResult) Extract() (*Rule, error) {
	var s struct {
		Rule *Rule `json:"firewall_rule"`
	}
	err := r.ExtractInto(&s)
	return s.Rule, err
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
