// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGroupCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := groups.CreateOpts{
		Name:     "testgroup",
		DomainID: "default",
		Extra: map[string]interface{}{
			"email": "testgroup@example.com",
		},
	}

	// Create Group in the default domain
	group, err := CreateGroup(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	tools.PrintResource(t, group)
	tools.PrintResource(t, group.Extra)

	updateOpts := groups.UpdateOpts{
		Description: "Test Users",
		Extra: map[string]interface{}{
			"email": "thetestgroup@example.com",
		},
	}

	newGroup, err := groups.Update(client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newGroup)
	tools.PrintResource(t, newGroup.Extra)

	listOpts := groups.ListOpts{
		DomainID: "default",
	}

	// List all Groups in default domain
	allPages, err := groups.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allGroups, err := groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	for _, g := range allGroups {
		tools.PrintResource(t, g)
		tools.PrintResource(t, g.Extra)
	}

	// Get the recently created group by ID
	p, err := groups.Get(client, group.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, p)
}
