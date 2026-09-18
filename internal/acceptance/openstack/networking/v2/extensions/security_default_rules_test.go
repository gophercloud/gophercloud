//go:build acceptance || networking || security

package extensions

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/defaultrules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestDefaultSecurityGroupRulesCreateGetListDelete(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	networking.RequireNeutronExtension(t, client, "security-groups-default-rules")

	description := "Default rule description"
	fromPort := tools.RandomInt(10080, 10089)
	toPort := tools.RandomInt(10090, 10099)

	createOpts := defaultrules.CreateOpts{
		Description:        description,
		Direction:          rules.DirIngress,
		EtherType:          rules.EtherType4,
		PortRangeMin:       ptr.To(fromPort),
		PortRangeMax:       ptr.To(toPort),
		Protocol:           rules.ProtocolTCP,
		RemoteGroupID:      "PARENT",
		UsedInDefaultSG:    ptr.To(true),
		UsedInNonDefaultSG: ptr.To(false),
	}

	rule, err := defaultrules.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err := defaultrules.Delete(context.TODO(), client, rule.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	tools.PrintResource(t, rule)

	th.AssertEquals(t, description, rule.Description)
	th.AssertEquals(t, "ingress", rule.Direction)
	th.AssertEquals(t, "IPv4", rule.EtherType)
	th.AssertEquals(t, fromPort, rule.PortRangeMin)
	th.AssertEquals(t, toPort, rule.PortRangeMax)
	th.AssertEquals(t, string(rules.ProtocolTCP), rule.Protocol)
	th.AssertEquals(t, "PARENT", rule.RemoteGroupID)
	th.AssertEquals(t, true, rule.UsedInDefaultSG)
	th.AssertEquals(t, false, rule.UsedInNonDefaultSG)

	getRule, err := defaultrules.Get(context.TODO(), client, rule.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, rule.ID, getRule.ID)
	th.AssertEquals(t, description, getRule.Description)

	listOpts := defaultrules.ListOpts{
		Protocol: string(rules.ProtocolTCP),
	}
	allPages, err := defaultrules.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allDefaultRules, err := defaultrules.ExtractDefaultRules(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, r := range allDefaultRules {
		if r.ID == rule.ID {
			found = true
			break
		}
	}
	th.AssertTrue(t, found)
}
