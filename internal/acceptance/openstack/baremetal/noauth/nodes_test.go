//go:build acceptance || baremetal || noauth

package noauth

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v1 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/baremetal/v1"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNodesCreateDestroy(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireIronicNoAuth(t)

	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.50"

	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)

	found := false
	err = nodes.List(client, nodes.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		nodeList, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		for _, n := range nodeList {
			if n.UUID == node.UUID {
				found = true
				return true, nil
			}
		}

		return false, nil
	})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, found, true)
}

func TestNodesUpdate(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireIronicNoAuth(t)

	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.50"

	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)

	updated, err := nodes.Update(context.TODO(), client, node.UUID, nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/maintenance",
			Value: "true",
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.Maintenance, true)
}

func TestNodesRAIDConfig(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ussuri")
	clients.RequireLong(t)
	clients.RequireIronicNoAuth(t)

	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.50"

	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)

	sizeGB := 100
	isTrue := true

	err = nodes.SetRAIDConfig(context.TODO(), client, node.UUID, nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				SizeGB:                &sizeGB,
				IsRootVolume:          &isTrue,
				RAIDLevel:             nodes.RAID5,
				DiskType:              nodes.HDD,
				NumberOfPhysicalDisks: 5,
			},
		},
	}).ExtractErr()
	th.AssertNoErr(t, err)
}
