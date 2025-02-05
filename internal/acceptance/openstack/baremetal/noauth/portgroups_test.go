package noauth

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v1 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/baremetal/v1"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/portgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPortGroupsCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)

	// NOTE(sharpz7) - increased due to create fake node requiring it.
	client.Microversion = "1.50"

	node, err := v1.CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)

	portgroup, err := v1.CreatePortGroup(t, client, node)
	th.AssertNoErr(t, err)
	defer v1.DeletePortGroup(t, client, portgroup)

	// Verify the portgroup exists by listing
	var found *portgroups.PortGroup
	err = portgroups.List(client, portgroups.ListOpts{
		Node: "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pg, err := portgroups.ExtractPortGroups(page)
		if err != nil {
			return false, err
		}

		for _, p := range pg {
			if p.UUID == portgroup.UUID {
				found = &p
				return true, nil
			}
		}

		return false, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, found, true)
}
