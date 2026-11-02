//go:build acceptance || baremetal || ports

package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPortsCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.53"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)
	port, err := CreatePort(t, client, node)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port)

	found := false
	err = ports.List(client, ports.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		portList, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}

		for _, p := range portList {
			if p.UUID == port.UUID {
				found = true
				return true, nil
			}
		}

		return false, nil
	})
	th.AssertNoErr(t, err)

	th.AssertTrue(t, found)
}

func TestPortsUpdate(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.53"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)
	port, err := CreatePort(t, client, node)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port)

	updated, err := ports.Update(context.TODO(), client, port.UUID, ports.UpdateOpts{
		ports.UpdateOperation{
			Op:    ports.ReplaceOp,
			Path:  "/address",
			Value: "aa:bb:cc:dd:ee:ff",
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "aa:bb:cc:dd:ee:ff", updated.Address)
}

func TestPortsPagination(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.53"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	// Cleanups run last-in first-out, so the node goes after its ports.
	t.Cleanup(func() { DeleteNode(t, client, node) })

	var created []string
	for i := range 3 {
		port, err := ports.Create(context.TODO(), client, ports.CreateOpts{
			NodeUUID: node.UUID,
			Address:  fmt.Sprintf("e6:72:1f:52:01:%02x", i),
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() { DeletePort(t, client, port) })
		created = append(created, port.UUID)
	}

	CheckPagination(t, func(limit int) pagination.Pager {
		return ports.List(client, ports.ListOpts{Limit: limit})
	}, func(page pagination.Page) ([]string, error) {
		list, err := ports.ExtractPorts(page)
		var ids []string
		for _, p := range list {
			ids = append(ids, p.UUID)
		}
		return ids, err
	}, created)
}
