//go:build acceptance || baremetal || nodes
// +build acceptance baremetal nodes

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestNodesCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	found := false
	err = nodes.List(client, nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
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

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	updated, err := nodes.Update(client, node.UUID, nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/maintenance",
			Value: "true",
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.Maintenance, true)
}

func TestNodesMaintenance(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	err = nodes.SetMaintenance(client, node.UUID, nodes.MaintenanceOpts{
		Reason: "I'm tired",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	updated, err := nodes.Get(client, node.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.Maintenance, true)
	th.AssertEquals(t, updated.MaintenanceReason, "I'm tired")

	err = nodes.UnsetMaintenance(client, node.UUID).ExtractErr()
	th.AssertNoErr(t, err)

	updated, err = nodes.Get(client, node.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.Maintenance, false)
	th.AssertEquals(t, updated.MaintenanceReason, "")
}

func TestNodesRAIDConfig(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ussuri")
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.50"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	sizeGB := 100
	isTrue := true

	err = nodes.SetRAIDConfig(client, node.UUID, nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				SizeGB:       &sizeGB,
				IsRootVolume: &isTrue,
				RAIDLevel:    nodes.RAID5,
				Controller:   "software",
				PhysicalDisks: []interface{}{
					map[string]string{
						"size": "> 100",
					},
					map[string]string{
						"size": "> 100",
					},
				},
			},
		},
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = nodes.SetRAIDConfig(client, node.UUID, nodes.RAIDConfigOpts{
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

func TestNodesFirmwareInterface(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/2023.2")
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.86"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	th.AssertEquals(t, node.FirmwareInterface, "no-firmware")

	nodeFirmwareCmps, err := nodes.ListFirmware(client, node.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, nodeFirmwareCmps, []nodes.FirmwareComponent{})
}
