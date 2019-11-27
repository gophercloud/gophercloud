// +build acceptance networking fwaas_v2

package fwaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestRuleCRUD(t *testing.T) {

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	rule, err := CreateRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	ruleDescription := "Some rule description"
	ruleProtocol := rules.ProtocolICMP
	updateOpts := rules.UpdateOpts{
		Description: &ruleDescription,
		Protocol:    &ruleProtocol,
	}

	ruleUpdated, err := rules.Update(client, rule.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, ruleUpdated.Description, ruleDescription)
	th.AssertEquals(t, ruleUpdated.Protocol, string(ruleProtocol))

	newRule, err := rules.Get(client, rule.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRule)
}
