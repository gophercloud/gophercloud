package rules

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a BandwidthLimitRule.
func (r commonResult) ExtractBandwidthLimitRule() (*BandwidthLimitRule, error) {
	var s struct {
		BandwidthLimitRule *BandwidthLimitRule `json:"bandwidth_limit_rule"`
	}
	err := r.ExtractInto(&s)
	return s.BandwidthLimitRule, err
}

// GetBandwidthLimitRuleResult represents the result of a Get operation. Call its Extract
// method to interpret it as a BandwidthLimitRule.
type GetBandwidthLimitRuleResult struct {
	commonResult
}

// CreateBandwidthLimitRuleResult represents the result of a Create operation. Call its Extract
// method to interpret it as a BandwidthLimitRule.
type CreateBandwidthLimitRuleResult struct {
	commonResult
}

// UpdateBandwidthLimitRuleResult represents the result of a Update operation. Call its Extract
// method to interpret it as a BandwidthLimitRule.
type UpdateBandwidthLimitRuleResult struct {
	commonResult
}

// DeleteBandwidthLimitRuleResult represents the result of a Delete operation. Call its Extract
// method to interpret it as a BandwidthLimitRule.
type DeleteBandwidthLimitRuleResult struct {
	gophercloud.ErrResult
}

// BandwidthLimitRule represents a QoS policy rule to set bandwidth limits.
type BandwidthLimitRule struct {
	// ID is a unique ID of the policy.
	ID string `json:"id"`

	// TenantID is the ID of the Identity project.
	TenantID string `json:"tenant_id"`

	// MaxKBps is a maximum kilobits per second.
	MaxKBps int `json:"max_kbps"`

	// MaxBurstKBps is a maximum burst size in kilobits.
	MaxBurstKBps int `json:"max_burst_kbps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction"`

	// Tags optionally set via extensions/attributestags.
	Tags []string `json:"tags"`
}

// BandwidthLimitRulePage stores a single page of BandwidthLimitRules from a List() API call.
type BandwidthLimitRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a BandwidthLimitRulePage is empty.
func (r BandwidthLimitRulePage) IsEmpty() (bool, error) {
	is, err := ExtractBandwidthLimitRules(r)
	return len(is) == 0, err
}

// ExtractBandwidthLimitRules accepts a BandwidthLimitRulePage, and extracts the elements into a slice of
// BandwidthLimitRules.
func ExtractBandwidthLimitRules(r pagination.Page) ([]BandwidthLimitRule, error) {
	var s []BandwidthLimitRule
	err := ExtractBandwidthLimitRulesInto(r, &s)
	return s, err
}

// ExtractBandwidthLimitRulesInto extracts the elements into a slice of RBAC Policy structs.
func ExtractBandwidthLimitRulesInto(r pagination.Page, v interface{}) error {
	return r.(BandwidthLimitRulePage).Result.ExtractIntoSlicePtr(v, "bandwidth_limit_rules")
}
