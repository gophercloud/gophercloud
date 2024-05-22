//go:build acceptance || baremetal || nodes

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNodesCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

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

	th.AssertEquals(t, node.ProvisionState, string(nodes.Enroll))

	err = nodes.ChangeProvisionState(context.TODO(), client, node.UUID, nodes.ProvisionStateOpts{
		Target: nodes.TargetManage,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = nodes.WaitForProvisionState(ctx, client, node.UUID, nodes.Manageable)
	th.AssertNoErr(t, err)
}

func TestNodesFields(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)
	err = nodes.List(client, nodes.ListOpts{
		Fields: []string{"uuid", "deploy_interface"},
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		nodeList, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		for _, n := range nodeList {
			if n.UUID == "" || n.DeployInterface == "" {
				t.Errorf("UUID or DeployInterface empty on %+v", n)
			}
			if n.BootInterface != "" {
				t.Errorf("BootInterface was not fetched but is not empty on %+v", n)
			}
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestNodesUpdate(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

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

func TestNodesMaintenance(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	err = nodes.SetMaintenance(context.TODO(), client, node.UUID, nodes.MaintenanceOpts{
		Reason: "I'm tired",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	updated, err := nodes.Get(context.TODO(), client, node.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, updated.Maintenance, true)
	th.AssertEquals(t, updated.MaintenanceReason, "I'm tired")

	err = nodes.UnsetMaintenance(context.TODO(), client, node.UUID).ExtractErr()
	th.AssertNoErr(t, err)

	updated, err = nodes.Get(context.TODO(), client, node.UUID).Extract()
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

	err = nodes.SetRAIDConfig(context.TODO(), client, node.UUID, nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				SizeGB:       &sizeGB,
				IsRootVolume: &isTrue,
				RAIDLevel:    nodes.RAID5,
				Controller:   "software",
				PhysicalDisks: []any{
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

	nodeFirmwareCmps, err := nodes.ListFirmware(context.TODO(), client, node.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, nodeFirmwareCmps, []nodes.FirmwareComponent{})
}

func TestNodesVirtualMedia(t *testing.T) {
	clients.SkipReleasesBelow(t, "master") // 2024.1
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.89"

	node, err := CreateNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	err = nodes.AttachVirtualMedia(context.TODO(), client, node.UUID, nodes.AttachVirtualMediaOpts{
		DeviceType: nodes.VirtualMediaCD,
		// It does not matter if QOTD server is actually present: the
		// request is processes asynchronously, all we need is a valid URL
		// that will not result in Ironic stuck for a long time.
		ImageURL: "http://127.0.0.1:17",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = nodes.DetachVirtualMedia(context.TODO(), client, node.UUID, nodes.DetachVirtualMediaOpts{
		DeviceTypes: []nodes.VirtualMediaDeviceType{nodes.VirtualMediaCD},
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = nodes.DetachVirtualMedia(context.TODO(), client, node.UUID, nodes.DetachVirtualMediaOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestNodesServicingHold(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/2023.2")
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.87"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	node, err = DeployFakeNode(t, client, node)
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	_, err = ChangeProvisionStateAndWait(ctx, client, node, nodes.ProvisionStateOpts{
		Target: nodes.TargetService,
		ServiceSteps: []nodes.ServiceStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      nodes.StepReboot,
			},
		},
	}, nodes.Active)
	th.AssertNoErr(t, err)
}
