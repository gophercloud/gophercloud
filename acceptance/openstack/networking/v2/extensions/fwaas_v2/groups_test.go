// +build acceptance networking fwaas_v2

package fwaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGroupCRUD(t *testing.T) {

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	createdGroup, err := CreateGroup(t, client)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, createdGroup.ID)

	tools.PrintResource(t, createdGroup)

	groupName := tools.RandomString("TESTACC-", 8)
	adminStateUp := false
	description := ("Some firewall group description")
	updateOpts := groups.UpdateOpts{
		Name:         &groupName,
		Description:  &description,
		AdminStateUp: &adminStateUp,
	}

	groupUpdated, err := groups.Update(client, createdGroup.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group %s: %v", createdGroup.ID, err)
	}

	th.AssertNoErr(t, err)
	th.AssertEquals(t, groupUpdated.Name, groupName)
	th.AssertEquals(t, groupUpdated.Description, description)
	th.AssertEquals(t, groupUpdated.AdminStateUp, adminStateUp)

	t.Logf("Updated firewall group %s", groupUpdated.ID)

	allPages, err := groups.List(client, nil).AllPages()
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
