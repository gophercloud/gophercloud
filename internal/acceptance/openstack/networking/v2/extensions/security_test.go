//go:build acceptance || networking || security

package extensions

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSecurityGroupsCreateUpdateDelete(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	group, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, group.ID)
	th.AssertEquals(t, group.Stateful, true)

	rule, err := CreateSecurityGroupRule(t, client, group.ID)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroupRule(t, client, rule.ID)

	rules, err := CreateSecurityGroupRulesBulk(t, client, group.ID)
	th.AssertNoErr(t, err)
	for _, r := range rules {
		defer DeleteSecurityGroupRule(t, client, r.ID)
	}

	tools.PrintResource(t, group)

	var name = "Update group"
	var description = ""
	updateOpts := groups.UpdateOpts{
		Name:        name,
		Description: &description,
		Stateful:    new(bool),
	}

	newGroup, err := groups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newGroup)
	th.AssertEquals(t, newGroup.Name, name)
	th.AssertEquals(t, newGroup.Description, description)
	th.AssertEquals(t, newGroup.Stateful, false)

	listOpts := groups.ListOpts{}
	allPages, err := groups.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err := groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, group := range allGroups {
		if group.ID == newGroup.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestSecurityGroupsPort(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	group, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, group.ID)

	rule, err := CreateSecurityGroupRule(t, client, group.ID)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroupRule(t, client, rule.ID)

	port, err := CreatePortWithSecurityGroup(t, client, network.ID, subnet.ID, group.ID)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)
}

func TestSecurityGroupsRevision(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a group
	group, err := CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSecurityGroup(t, client, group.ID)

	tools.PrintResource(t, group)

	// Store the current revision number.
	oldRevisionNumber := group.RevisionNumber

	// Update the group without revision number.
	// This should work.
	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &groups.UpdateOpts{
		Name:        newName,
		Description: &newDescription,
	}
	group, err = groups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, group)

	// This should fail due to an old revision number.
	newDescription = "new description"
	updateOpts = &groups.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &oldRevisionNumber,
	}
	_, err = groups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertErr(t, err)
	if !strings.Contains(err.Error(), "RevisionNumberConstraintFailed") {
		t.Fatalf("expected to see an error of type RevisionNumberConstraintFailed, but got the following error instead: %v", err)
	}

	// Reread the group to show that it did not change.
	group, err = groups.Get(context.TODO(), client, group.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, group)

	// This should work because now we do provide a valid revision number.
	newDescription = "new description"
	updateOpts = &groups.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &group.RevisionNumber,
	}
	group, err = groups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, group)

	th.AssertEquals(t, group.Name, newName)
	th.AssertEquals(t, group.Description, newDescription)
}
