//go:build acceptance || baremetal || portgroups

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
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

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	portgroup, err := CreatePortGroup(t, client, node)
	th.AssertNoErr(t, err)
	defer DeletePortGroup(t, client, portgroup)

	// Verify the portgroup exists by listing
	err = portgroups.List(client, portgroups.ListOpts{
		Node: node.UUID,
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pg, err := portgroups.ExtractPortGroups(page)
		if err != nil {
			return false, err
		}

		for _, p := range pg {
			if p.UUID == portgroup.UUID {
				return true, nil
			}
		}

		return false, nil
	})
	th.AssertNoErr(t, err)
}

func TestPortGroupsPagination(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.50"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	// Cleanups run last-in first-out, so the node goes after its ports.
	t.Cleanup(func() { DeleteNode(t, client, node) })

	var created []string
	for range 3 {
		portgroup, err := CreatePortGroup(t, client, node)
		th.AssertNoErr(t, err)
		t.Cleanup(func() { DeletePortGroup(t, client, portgroup) })
		created = append(created, portgroup.UUID)
	}

	CheckPagination(t, func(limit int) pagination.Pager {
		return portgroups.List(client, portgroups.ListOpts{Limit: limit})
	}, func(page pagination.Page) ([]string, error) {
		list, err := portgroups.ExtractPortGroups(page)
		var ids []string
		for _, pg := range list {
			ids = append(ids, pg.UUID)
		}
		return ids, err
	}, created)
}
