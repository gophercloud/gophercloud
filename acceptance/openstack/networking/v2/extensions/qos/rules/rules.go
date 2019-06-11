package rules

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
	th "github.com/gophercloud/gophercloud/testhelper"
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

	rule, err := rules.CreateBandwidthLimitRule(client, policyID, createOpts).ExtractBandwidthLimitRule()
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

	rule, err := rules.CreateDSCPMarkingRule(client, policyID, createOpts).ExtractDSCPMarkingRule()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS DSCP marking rule")

	th.AssertEquals(t, dscpMark, rule.DSCPMark)

	return rule, nil
}
