//go:build acceptance || networking || fwaas_v2

package fwaas_v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/groups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGroupCRUD(t *testing.T) {
	// Releases below Victoria are not maintained.
	// FWaaS_v2 is not compatible with releases below Zed.
	clients.SkipReleasesBelow(t, "stable/zed")

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "fwaas_v2")

	createdGroup, err := CreateGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, createdGroup.ID)

	tools.PrintResource(t, createdGroup)

	createdRule, err := CreateRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRule(t, client, createdRule.ID)

	tools.PrintResource(t, createdRule)

	createdPolicy, err := CreatePolicy(t, client, createdRule.ID)
	th.AssertNoErr(t, err)
	defer DeletePolicy(t, client, createdPolicy.ID)

	tools.PrintResource(t, createdPolicy)

	groupName := tools.RandomString("TESTACC-", 8)
	adminStateUp := false
	description := ("Some firewall group description")
	firewall_policy_id := createdPolicy.ID
	updateOpts := groups.UpdateOpts{
		Name:                    &groupName,
		Description:             &description,
		AdminStateUp:            &adminStateUp,
		IngressFirewallPolicyID: &firewall_policy_id,
		EgressFirewallPolicyID:  &firewall_policy_id,
	}

	updatedGroup, err := groups.Update(context.TODO(), client, createdGroup.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group %s: %v", createdGroup.ID, err)
	}

	th.AssertNoErr(t, err)
	th.AssertEquals(t, updatedGroup.Name, groupName)
	th.AssertEquals(t, updatedGroup.Description, description)
	th.AssertEquals(t, updatedGroup.AdminStateUp, adminStateUp)
	th.AssertEquals(t, updatedGroup.IngressFirewallPolicyID, firewall_policy_id)
	th.AssertEquals(t, updatedGroup.EgressFirewallPolicyID, firewall_policy_id)

	t.Logf("Updated firewall group %s", updatedGroup.ID)

	removeIngressPolicy, err := groups.RemoveIngressPolicy(context.TODO(), client, updatedGroup.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to remove ingress firewall policy from firewall group %s: %v", removeIngressPolicy.ID, err)
	}

	th.AssertEquals(t, removeIngressPolicy.IngressFirewallPolicyID, "")
	th.AssertEquals(t, removeIngressPolicy.EgressFirewallPolicyID, firewall_policy_id)

	t.Logf("Ingress policy removed from firewall group %s", updatedGroup.ID)

	removeEgressPolicy, err := groups.RemoveEgressPolicy(context.TODO(), client, updatedGroup.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to remove egress firewall policy from firewall group %s: %v", removeEgressPolicy.ID, err)
	}

	th.AssertEquals(t, removeEgressPolicy.EgressFirewallPolicyID, "")

	t.Logf("Egress policy removed from firewall group %s", updatedGroup.ID)

	allPages, err := groups.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err := groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Attempting to find firewall group %s\n", createdGroup.ID)
	var found bool
	for _, group := range allGroups {
		if group.ID == createdGroup.ID {
			found = true
			t.Logf("Found firewall group %s\n", group.ID)
		}
	}

	th.AssertEquals(t, found, true)
}
