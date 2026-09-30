package rules

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateBandwidthLimitRule will create a QoS BandwidthLimitRule associated with the provided QoS policy.
// An error will be returned if the QoS rule could not be created.
func CreateBandwidthLimitRule(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*rules.BandwidthLimitRule, error) {
	maxKBps := 3000
	maxBurstKBps := 300

	createOpts := rules.CreateBandwidthLimitRuleOpts{
		MaxKBps:      maxKBps,
		MaxBurstKBps: maxBurstKBps,
	}

	t.Logf("Attempting to create a QoS bandwidth limit rule with max_kbps: %d, max_burst_kbps: %d", maxKBps, maxBurstKBps)

	rule, err := rules.CreateBandwidthLimitRule(context.TODO(), client, policyID, createOpts).ExtractBandwidthLimitRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS bandwidth limit rule")

	th.AssertEquals(t, maxKBps, rule.MaxKBps)
	th.AssertEquals(t, maxBurstKBps, rule.MaxBurstKBps)

	return rule, nil
}

// CreateDSCPMarkingRule will create a QoS DSCPMarkingRule associated with the provided QoS policy.
// An error will be returned if the QoS rule could not be created.
func CreateDSCPMarkingRule(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*rules.DSCPMarkingRule, error) {
	dscpMark := 26

	createOpts := rules.CreateDSCPMarkingRuleOpts{
		DSCPMark: dscpMark,
	}

	t.Logf("Attempting to create a QoS DSCP marking rule with dscp_mark: %d", dscpMark)

	rule, err := rules.CreateDSCPMarkingRule(context.TODO(), client, policyID, createOpts).ExtractDSCPMarkingRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS DSCP marking rule")

	th.AssertEquals(t, dscpMark, rule.DSCPMark)

	return rule, nil
}

// CreateMinimumBandwidthRule will create a QoS MinimumBandwidthRule associated with the provided QoS policy.
// An error will be returned if the QoS rule could not be created.
func CreateMinimumBandwidthRule(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*rules.MinimumBandwidthRule, error) {
	minKBps := 1000

	createOpts := rules.CreateMinimumBandwidthRuleOpts{
		MinKBps: minKBps,
	}

	t.Logf("Attempting to create a QoS minimum bandwidth rule with min_kbps: %d", minKBps)

	rule, err := rules.CreateMinimumBandwidthRule(context.TODO(), client, policyID, createOpts).ExtractMinimumBandwidthRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS minimum bandwidth rule")

	th.AssertEquals(t, minKBps, rule.MinKBps)

	return rule, nil
}

// CreateMinimumPacketRateRule will create a QoS MinimumPacketRateRule associated with the provided QoS policy.
// An error will be returned if the QoS rule could not be created.
func CreateMinimumPacketRateRule(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*rules.MinimumPacketRateRule, error) {
	minKPps := 1000
	// The minimum packet rate rule supports the "any" direction, which is
	// specific to this rule type.
	direction := "any"

	createOpts := rules.CreateMinimumPacketRateRuleOpts{
		MinKPps:   minKPps,
		Direction: direction,
	}

	t.Logf("Attempting to create a QoS minimum packet rate rule with min_kpps: %d, direction: %s", minKPps, direction)

	rule, err := rules.CreateMinimumPacketRateRule(context.TODO(), client, policyID, createOpts).ExtractMinimumPacketRateRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS minimum packet rate rule")

	th.AssertEquals(t, minKPps, rule.MinKPps)
	th.AssertEquals(t, direction, rule.Direction)

	return rule, nil
}

// CreatePacketRateLimitRule will create a QoS PacketRateLimitRule associated with the provided QoS policy.
// An error will be returned if the QoS rule could not be created.
func CreatePacketRateLimitRule(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*rules.PacketRateLimitRule, error) {
	maxKPps := 3000
	maxBurstKPps := 300

	createOpts := rules.CreatePacketRateLimitRuleOpts{
		MaxKPps:      maxKPps,
		MaxBurstKPps: maxBurstKPps,
	}

	t.Logf("Attempting to create a QoS packet rate limit rule with max_kpps: %d, max_burst_kpps: %d", maxKPps, maxBurstKPps)

	rule, err := rules.CreatePacketRateLimitRule(context.TODO(), client, policyID, createOpts).ExtractPacketRateLimitRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS packet rate limit rule")

	th.AssertEquals(t, maxKPps, rule.MaxKPps)
	th.AssertEquals(t, maxBurstKPps, rule.MaxBurstKPps)

	return rule, nil
}
