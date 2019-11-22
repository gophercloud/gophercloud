package rules

import (
	"github.com/gophercloud/gophercloud"
)

// Rule represents a firewall rule
type Rule struct {
	ID                   string `json:"id"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	Protocol             string `json:"protocol"`
	Action               string `json:"action"`
	IPVersion            int    `json:"ip_version,omitempty"`
	SourceIPAddress      string `json:"source_ip_address,omitempty"`
	DestinationIPAddress string `json:"destination_ip_address,omitempty"`
	SourcePort           string `json:"source_port,omitempty"`
	DestinationPort      string `json:"destination_port,omitempty"`
	Shared               bool   `json:"shared,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	FirewallPolicyID     string `json:"firewall_policy_id"`
	TenantID             string `json:"tenant_id"`
	ProjectID            string `json:"project_id"`
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

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}
