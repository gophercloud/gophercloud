package rules

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	accpolicies "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestBandwidthLimitRulesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a QoS policy
	policy, err := accpolicies.CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer policies.Delete(client, policy.ID)

	tools.PrintResource(t, policy)

	// Create a QoS policy rule.
	rule, err := CreateBandwidthLimitRule(t, client, policy.ID)
	th.AssertNoErr(t, err)
	defer rules.DeleteBandwidthLimitRule(client, policy.ID, rule.ID)

	// Update the QoS policy rule.
	newMaxBurstKBps := 0
	updateOpts := rules.UpdateBandwidthLimitRuleOpts{
		MaxBurstKBps: &newMaxBurstKBps,
	}
	newRule, err := rules.UpdateBandwidthLimitRule(client, policy.ID, rule.ID, updateOpts).ExtractBandwidthLimitRule()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)
	th.AssertEquals(t, newRule.MaxBurstKBps, 0)

	allPages, err := rules.ListBandwidthLimitRules(client, policy.ID, rules.BandwidthLimitRulesListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allRules, err := rules.ExtractBandwidthLimitRules(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, rule := range allRules {
		if rule.ID == newRule.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
