// +build acceptance compute secgroups

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestSecGroupsList(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := secgroups.List(client).AllPages()
	th.AssertNoErr(t, err)

	allSecGroups, err := secgroups.ExtractSecurityGroups(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, secgroup := range allSecGroups {
		tools.PrintResource(t, secgroup)

		if secgroup.Name == "default" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestSecGroupsCRUD(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	securityGroup, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, securityGroup.ID)

	tools.PrintResource(t, securityGroup)

	newName := tools.RandomString("secgroup_", 4)
	updateOpts := secgroups.UpdateOpts{
		Name:        newName,
		Description: tools.RandomString("dec_", 10),
	}
	updatedSecurityGroup, err := secgroups.Update(client, securityGroup.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedSecurityGroup)

	t.Logf("Updated %s's name to %s", updatedSecurityGroup.ID, updatedSecurityGroup.Name)

	th.AssertEquals(t, updatedSecurityGroup.Name, newName)
}

func TestSecGroupsRuleCreate(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	securityGroup, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, securityGroup.ID)

	tools.PrintResource(t, securityGroup)

	rule, err := CreateSecurityGroupRule(t, client, securityGroup.ID)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroupRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	newSecurityGroup, err := secgroups.Get(client, securityGroup.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newSecurityGroup)

	th.AssertEquals(t, len(newSecurityGroup.Rules), 1)
}

func TestSecGroupsAddGroupToServer(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	securityGroup, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, securityGroup.ID)

	rule, err := CreateSecurityGroupRule(t, client, securityGroup.ID)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroupRule(t, client, rule.ID)

	server, err := CreateServer(t, client)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer DeleteServer(t, client, server)

	t.Logf("Adding group %s to server %s", securityGroup.ID, server.ID)
	err = secgroups.AddServer(client, server.ID, securityGroup.Name).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to add group %s to server %s: %s", securityGroup.ID, server.ID, err)
	}

	server, err = servers.Get(client, server.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get server %s: %s", server.ID, err)
	}

	tools.PrintResource(t, server)

	t.Logf("Removing group %s from server %s", securityGroup.ID, server.ID)
	err = secgroups.RemoveServer(client, server.ID, securityGroup.Name).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to remove group %s from server %s: %s", securityGroup.ID, server.ID, err)
	}
}
