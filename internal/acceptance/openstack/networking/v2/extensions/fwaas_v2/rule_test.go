//go:build acceptance || networking || fwaas_v2

package fwaas_v2

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/rules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRuleCRUD(t *testing.T) {
	// Releases below Victoria are not maintained.
	// FWaaS_v2 is not compatible with releases below Zed.
	clients.SkipReleasesBelow(t, "stable/zed")

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "fwaas_v2")

	rule, err := CreateRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	ruleDescription := "Some rule description"
	ruleSourcePortInt := strconv.Itoa(tools.RandomInt(1, 100))
	ruleSourcePort := fmt.Sprintf("%s:%s", ruleSourcePortInt, ruleSourcePortInt)
	ruleProtocol := rules.ProtocolTCP
	updateOpts := rules.UpdateOpts{
		Description: &ruleDescription,
		Protocol:    &ruleProtocol,
		SourcePort:  &ruleSourcePort,
	}

	ruleUpdated, err := rules.Update(context.TODO(), client, rule.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, ruleUpdated.Description, ruleDescription)
	th.AssertEquals(t, ruleUpdated.SourcePort, ruleSourcePortInt)
	th.AssertEquals(t, ruleUpdated.Protocol, string(ruleProtocol))

	newRule, err := rules.Get(context.TODO(), client, rule.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)

	allPages, err := rules.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRules, err := rules.ExtractRules(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Attempting to find rule %s\n", newRule.ID)
	var found bool
	for _, rule := range allRules {
		if rule.ID == newRule.ID {
			found = true
			t.Logf("Found rule %s\n", newRule.ID)
		}
	}

	th.AssertEquals(t, found, true)
}
