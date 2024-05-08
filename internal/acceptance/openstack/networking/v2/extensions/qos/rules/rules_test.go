//go:build acceptance || networking || qos || rules

package rules

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	accpolicies "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBandwidthLimitRulesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "qos")

	// Create a QoS policy
	policy, err := accpolicies.CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer policies.Delete(context.TODO(), client, policy.ID)

	tools.PrintResource(t, policy)

	// Create a QoS policy rule.
	rule, err := CreateBandwidthLimitRule(t, client, policy.ID)
	th.AssertNoErr(t, err)
	defer rules.DeleteBandwidthLimitRule(context.TODO(), client, policy.ID, rule.ID)

	// Update the QoS policy rule.
	newMaxBurstKBps := 0
	updateOpts := rules.UpdateBandwidthLimitRuleOpts{
		MaxBurstKBps: &newMaxBurstKBps,
	}
	newRule, err := rules.UpdateBandwidthLimitRule(context.TODO(), client, policy.ID, rule.ID, updateOpts).ExtractBandwidthLimitRule()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)
	th.AssertEquals(t, newRule.MaxBurstKBps, 0)

	allPages, err := rules.ListBandwidthLimitRules(client, policy.ID, rules.BandwidthLimitRulesListOpts{}).AllPages(context.TODO())
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

func TestDSCPMarkingRulesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "qos")

	// Create a QoS policy
	policy, err := accpolicies.CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer policies.Delete(context.TODO(), client, policy.ID)

	tools.PrintResource(t, policy)

	// Create a QoS policy rule.
	rule, err := CreateDSCPMarkingRule(t, client, policy.ID)
	th.AssertNoErr(t, err)
	defer rules.DeleteDSCPMarkingRule(context.TODO(), client, policy.ID, rule.ID)

	// Update the QoS policy rule.
	dscpMark := 20
	updateOpts := rules.UpdateDSCPMarkingRuleOpts{
		DSCPMark: &dscpMark,
	}
	newRule, err := rules.UpdateDSCPMarkingRule(context.TODO(), client, policy.ID, rule.ID, updateOpts).ExtractDSCPMarkingRule()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)
	th.AssertEquals(t, newRule.DSCPMark, 20)

	allPages, err := rules.ListDSCPMarkingRules(client, policy.ID, rules.DSCPMarkingRulesListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRules, err := rules.ExtractDSCPMarkingRules(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, rule := range allRules {
		if rule.ID == newRule.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestMinimumBandwidthRulesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "qos")

	// Create a QoS policy
	policy, err := accpolicies.CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer policies.Delete(context.TODO(), client, policy.ID)

	tools.PrintResource(t, policy)

	// Create a QoS policy rule.
	rule, err := CreateMinimumBandwidthRule(t, client, policy.ID)
	th.AssertNoErr(t, err)
	defer rules.DeleteMinimumBandwidthRule(context.TODO(), client, policy.ID, rule.ID)

	// Update the QoS policy rule.
	minKBps := 500
	updateOpts := rules.UpdateMinimumBandwidthRuleOpts{
		MinKBps: &minKBps,
	}
	newRule, err := rules.UpdateMinimumBandwidthRule(context.TODO(), client, policy.ID, rule.ID, updateOpts).ExtractMinimumBandwidthRule()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)
	th.AssertEquals(t, newRule.MinKBps, 500)

	allPages, err := rules.ListMinimumBandwidthRules(client, policy.ID, rules.MinimumBandwidthRulesListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRules, err := rules.ExtractMinimumBandwidthRules(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, rule := range allRules {
		if rule.ID == newRule.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
