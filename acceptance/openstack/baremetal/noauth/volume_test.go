package noauth

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	v1 "github.com/gophercloud/gophercloud/acceptance/openstack/baremetal/v1"
	"github.com/gophercloud/gophercloud/acceptance/openstack/blockstorage/v3"
	bmvolume "github.com/gophercloud/gophercloud/openstack/baremetal/v1/volume"
	"github.com/gophercloud/gophercloud/pagination"

	th "github.com/gophercloud/gophercloud/testhelper"
)

// node storage interface should be set to cinder or external

func TestConnectorCreateDestroy(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"
	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	unode, err := v1.UpdateNodeStorageInterface(client, node.UUID, "cinder")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "cinder", unode.StorageInterface)
	connector, err := v1.CreateVolumeConnector(t, client, node)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)
	defer v1.DeleteVolumeConnector(t, client, connector)
	found := false
	err = bmvolume.ListConnectors(client, bmvolume.ListConnectorsOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		connectorList, err := bmvolume.ExtractConnectors(page)
		if err != nil {
			return false, err
		}
		for _, c := range connectorList {
			if c.UUID == connector.UUID {
				found = true
				return true, nil
			}
		}
		return false, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, found, true)
}

func TestConnectorUpdate(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"
	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	unode, err := v1.UpdateNodeStorageInterface(client, node.UUID, "cinder")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "cinder", unode.StorageInterface)
	connector, err := v1.CreateVolumeConnector(t, client, node)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)
	defer v1.DeleteVolumeConnector(t, client, connector)
	updated, err := bmvolume.UpdateConnector(client, connector.UUID, bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/connector_id",
			Value: "iqn.2017-08.org.openstack." + node.UUID,
		},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "iqn.2017-08.org.openstack."+node.UUID, updated.ConnectorId)
}

func TestTargetCreateDestroy(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"
	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	unode, err := v1.UpdateNodeStorageInterface(client, node.UUID, "cinder")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "cinder", unode.StorageInterface)
	volume, err := v3.CreateVolume(t, client)
	th.AssertNoErr(t, err)
	target, err := v1.CreateVolumeTarget(t, client, node, volume.ID)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)
	defer v1.DeleteVolumeTarget(t, client, target)
	defer v3.DeleteVolume(t, client, volume)
	found := false
	err = bmvolume.ListTargets(client, bmvolume.ListTargetsOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		targetList, err := bmvolume.ExtractTargets(page)
		if err != nil {
			return false, err
		}
		for _, c := range targetList {
			if c.UUID == target.UUID {
				found = true
				return true, nil
			}
		}
		return false, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, found, true)
}

func TestTargetUpdate(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewBareMetalV1NoAuthClient()
	th.AssertNoErr(t, err)
	client.Microversion = "1.38"
	node, err := v1.CreateNode(t, client)
	th.AssertNoErr(t, err)
	unode, err := v1.UpdateNodeStorageInterface(client, node.UUID, "cinder")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "cinder", unode.StorageInterface)
	volume, err := v3.CreateVolume(t, client)
	th.AssertNoErr(t, err)
	target, err := v1.CreateVolumeTarget(t, client, node, volume.ID)
	th.AssertNoErr(t, err)
	defer v1.DeleteNode(t, client, node)
	defer v1.DeleteVolumeTarget(t, client, target)
	defer v3.DeleteVolume(t, client, volume)
	another_volume, err := v3.CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer v3.DeleteVolume(t, client, another_volume)
	updated, err := bmvolume.UpdateTarget(client, target.UUID, bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/volume_id",
			Value: another_volume.ID,
		},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, another_volume.ID, updated.VolumeId)
}
