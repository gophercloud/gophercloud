package rules

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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
	if r.StatusCode == 204 {
		return true, nil
	}

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
func ExtractBandwidthLimitRulesInto(r pagination.Page, v any) error {
	return r.(BandwidthLimitRulePage).ExtractIntoSlicePtr(v, "bandwidth_limit_rules")
}

// Extract is a function that accepts a result and extracts a DSCPMarkingRule.
func (r commonResult) ExtractDSCPMarkingRule() (*DSCPMarkingRule, error) {
	var s struct {
		DSCPMarkingRule *DSCPMarkingRule `json:"dscp_marking_rule"`
	}
	err := r.ExtractInto(&s)
	return s.DSCPMarkingRule, err
}

// GetDSCPMarkingRuleResult represents the result of a Get operation. Call its Extract
// method to interpret it as a DSCPMarkingRule.
type GetDSCPMarkingRuleResult struct {
	commonResult
}

// CreateDSCPMarkingRuleResult represents the result of a Create operation. Call its Extract
// method to interpret it as a DSCPMarkingRule.
type CreateDSCPMarkingRuleResult struct {
	commonResult
}

// UpdateDSCPMarkingRuleResult represents the result of a Update operation. Call its Extract
// method to interpret it as a DSCPMarkingRule.
type UpdateDSCPMarkingRuleResult struct {
	commonResult
}

// DeleteDSCPMarkingRuleResult represents the result of a Delete operation. Call its Extract
// method to interpret it as a DSCPMarkingRule.
type DeleteDSCPMarkingRuleResult struct {
	gophercloud.ErrResult
}

// DSCPMarkingRule represents a QoS policy rule to set DSCP marking.
type DSCPMarkingRule struct {
	// ID is a unique ID of the policy.
	ID string `json:"id"`

	// TenantID is the ID of the Identity project.
	TenantID string `json:"tenant_id"`

	// DSCPMark contains DSCP mark value.
	DSCPMark int `json:"dscp_mark"`

	// Tags optionally set via extensions/attributestags.
	Tags []string `json:"tags"`
}

// DSCPMarkingRulePage stores a single page of DSCPMarkingRules from a List() API call.
type DSCPMarkingRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a DSCPMarkingRulePage is empty.
func (r DSCPMarkingRulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractDSCPMarkingRules(r)
	return len(is) == 0, err
}

// ExtractDSCPMarkingRules accepts a DSCPMarkingRulePage, and extracts the elements into a slice of
// DSCPMarkingRules.
func ExtractDSCPMarkingRules(r pagination.Page) ([]DSCPMarkingRule, error) {
	var s []DSCPMarkingRule
	err := ExtractDSCPMarkingRulesInto(r, &s)
	return s, err
}

// ExtractDSCPMarkingRulesInto extracts the elements into a slice of RBAC Policy structs.
func ExtractDSCPMarkingRulesInto(r pagination.Page, v any) error {
	return r.(DSCPMarkingRulePage).ExtractIntoSlicePtr(v, "dscp_marking_rules")
}

// Extract is a function that accepts a result and extracts a BandwidthLimitRule.
func (r commonResult) ExtractMinimumBandwidthRule() (*MinimumBandwidthRule, error) {
	var s struct {
		MinimumBandwidthRule *MinimumBandwidthRule `json:"minimum_bandwidth_rule"`
	}
	err := r.ExtractInto(&s)
	return s.MinimumBandwidthRule, err
}

// GetMinimumBandwidthRuleResult represents the result of a Get operation. Call its Extract
// method to interpret it as a MinimumBandwidthRule.
type GetMinimumBandwidthRuleResult struct {
	commonResult
}

// CreateMinimumBandwidthRuleResult represents the result of a Create operation. Call its Extract
// method to interpret it as a MinimumBandwidthtRule.
type CreateMinimumBandwidthRuleResult struct {
	commonResult
}

// UpdateMinimumBandwidthRuleResult represents the result of a Update operation. Call its Extract
// method to interpret it as a MinimumBandwidthRule.
type UpdateMinimumBandwidthRuleResult struct {
	commonResult
}

// DeleteMinimumBandwidthRuleResult represents the result of a Delete operation. Call its Extract
// method to interpret it as a MinimumBandwidthRule.
type DeleteMinimumBandwidthRuleResult struct {
	gophercloud.ErrResult
}

// MinimumBandwidthRule represents a QoS policy rule to set minimum bandwidth.
type MinimumBandwidthRule struct {
	// ID is a unique ID of the rule.
	ID string `json:"id"`

	// TenantID is the ID of the Identity project.
	TenantID string `json:"tenant_id"`

	// MaxKBps is a maximum kilobits per second.
	MinKBps int `json:"min_kbps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction"`

	// Tags optionally set via extensions/attributestags.
	Tags []string `json:"tags"`
}

// MinimumBandwidthRulePage stores a single page of MinimumBandwidthRules from a List() API call.
type MinimumBandwidthRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a MinimumBandwidthRulePage is empty.
func (r MinimumBandwidthRulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractMinimumBandwidthRules(r)
	return len(is) == 0, err
}

// ExtractMinimumBandwidthRules accepts a MinimumBandwidthRulePage, and extracts the elements into a slice of
// MinimumBandwidthRules.
func ExtractMinimumBandwidthRules(r pagination.Page) ([]MinimumBandwidthRule, error) {
	var s []MinimumBandwidthRule
	err := ExtractMinimumBandwidthRulesInto(r, &s)
	return s, err
}

// ExtractMinimumBandwidthRulesInto extracts the elements into a slice of RBAC Policy structs.
func ExtractMinimumBandwidthRulesInto(r pagination.Page, v any) error {
	return r.(MinimumBandwidthRulePage).ExtractIntoSlicePtr(v, "minimum_bandwidth_rules")
}

// ExtractMinimumPacketRateRule is a function that accepts a result and extracts a MinimumPacketRateRule.
func (r commonResult) ExtractMinimumPacketRateRule() (*MinimumPacketRateRule, error) {
	var s struct {
		MinimumPacketRateRule *MinimumPacketRateRule `json:"minimum_packet_rate_rule"`
	}
	err := r.ExtractInto(&s)
	return s.MinimumPacketRateRule, err
}

// GetMinimumPacketRateRuleResult represents the result of a Get operation. Call its Extract
// method to interpret it as a MinimumPacketRateRule.
type GetMinimumPacketRateRuleResult struct {
	commonResult
}

// CreateMinimumPacketRateRuleResult represents the result of a Create operation. Call its Extract
// method to interpret it as a MinimumPacketRateRule.
type CreateMinimumPacketRateRuleResult struct {
	commonResult
}

// UpdateMinimumPacketRateRuleResult represents the result of a Update operation. Call its Extract
// method to interpret it as a MinimumPacketRateRule.
type UpdateMinimumPacketRateRuleResult struct {
	commonResult
}

// DeleteMinimumPacketRateRuleResult represents the result of a Delete operation. Call its Extract
// method to interpret it as a MinimumPacketRateRule.
type DeleteMinimumPacketRateRuleResult struct {
	gophercloud.ErrResult
}

// MinimumPacketRateRule represents a QoS policy rule to set minimum packet rate.
type MinimumPacketRateRule struct {
	// ID is a unique ID of the rule.
	ID string `json:"id"`

	// MinKPps is a minimum kilopackets per second.
	MinKPps int `json:"min_kpps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction"`

	// Tags optionally set via extensions/attributestags.
	Tags []string `json:"tags"`
}

// MinimumPacketRateRulePage stores a single page of MinimumPacketRateRules from a List() API call.
type MinimumPacketRateRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a MinimumPacketRateRulePage is empty.
func (r MinimumPacketRateRulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractMinimumPacketRateRules(r)
	return len(is) == 0, err
}

// ExtractMinimumPacketRateRules accepts a MinimumPacketRateRulePage, and extracts the elements into a slice of
// MinimumPacketRateRules.
func ExtractMinimumPacketRateRules(r pagination.Page) ([]MinimumPacketRateRule, error) {
	var s []MinimumPacketRateRule
	err := ExtractMinimumPacketRateRulesInto(r, &s)
	return s, err
}

// ExtractMinimumPacketRateRulesInto extracts the elements into a slice of MinimumPacketRateRule structs.
func ExtractMinimumPacketRateRulesInto(r pagination.Page, v any) error {
	return r.(MinimumPacketRateRulePage).ExtractIntoSlicePtr(v, "minimum_packet_rate_rules")
}

// ExtractPacketRateLimitRule is a function that accepts a result and extracts a PacketRateLimitRule.
func (r commonResult) ExtractPacketRateLimitRule() (*PacketRateLimitRule, error) {
	var s struct {
		PacketRateLimitRule *PacketRateLimitRule `json:"packet_rate_limit_rule"`
	}
	err := r.ExtractInto(&s)
	return s.PacketRateLimitRule, err
}

// GetPacketRateLimitRuleResult represents the result of a Get operation. Call its Extract
// method to interpret it as a PacketRateLimitRule.
type GetPacketRateLimitRuleResult struct {
	commonResult
}

// CreatePacketRateLimitRuleResult represents the result of a Create operation. Call its Extract
// method to interpret it as a PacketRateLimitRule.
type CreatePacketRateLimitRuleResult struct {
	commonResult
}

// UpdatePacketRateLimitRuleResult represents the result of a Update operation. Call its Extract
// method to interpret it as a PacketRateLimitRule.
type UpdatePacketRateLimitRuleResult struct {
	commonResult
}

// DeletePacketRateLimitRuleResult represents the result of a Delete operation. Call its Extract
// method to interpret it as a PacketRateLimitRule.
type DeletePacketRateLimitRuleResult struct {
	gophercloud.ErrResult
}

// PacketRateLimitRule represents a QoS policy rule to set packet rate limits.
type PacketRateLimitRule struct {
	// ID is a unique ID of the rule.
	ID string `json:"id"`

	// MaxKPps is a maximum kilopackets per second.
	MaxKPps int `json:"max_kpps"`

	// MaxBurstKPps is a maximum burst size in kilopackets.
	MaxBurstKPps int `json:"max_burst_kpps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction"`

	// Tags optionally set via extensions/attributestags.
	Tags []string `json:"tags"`
}

// PacketRateLimitRulePage stores a single page of PacketRateLimitRules from a List() API call.
type PacketRateLimitRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a PacketRateLimitRulePage is empty.
func (r PacketRateLimitRulePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractPacketRateLimitRules(r)
	return len(is) == 0, err
}

// ExtractPacketRateLimitRules accepts a PacketRateLimitRulePage, and extracts the elements into a slice of
// PacketRateLimitRules.
func ExtractPacketRateLimitRules(r pagination.Page) ([]PacketRateLimitRule, error) {
	var s []PacketRateLimitRule
	err := ExtractPacketRateLimitRulesInto(r, &s)
	return s, err
}

// ExtractPacketRateLimitRulesInto extracts the elements into a slice of PacketRateLimitRule structs.
func ExtractPacketRateLimitRulesInto(r pagination.Page, v any) error {
	return r.(PacketRateLimitRulePage).ExtractIntoSlicePtr(v, "packet_rate_limit_rules")
}
