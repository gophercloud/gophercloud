//go:build acceptance || compute || volumeattach

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	bs "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/blockstorage/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumeAttachAttachment(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	// TEMP: Check inventory to figure out what's missing
	placementClient, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	placementClient.Microversion = "1.37"

	allPages, err := resourceproviders.List(placementClient, resourceproviders.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	th.AssertNoErr(t, err)

	for _, r := range allResourceProviders {
		inv, err := resourceproviders.GetInventories(context.TODO(), placementClient, r.UUID).Extract()
		th.AssertNoErr(t, err)

		for k, v := range inv.Inventories {
			t.Logf("resource %s (%s) inventory: %s=%+v", r.UUID, r.Name, k, v)
		}

		use, err := resourceproviders.GetUsages(context.TODO(), placementClient, r.UUID).Extract()
		th.AssertNoErr(t, err)

		for k, v := range use.Usages {
			t.Logf("resource %s (%s) usage: %s=%d", r.UUID, r.Name, k, v)
		}
	}
	// TEMP: Check inventory to figure out what's missing

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	volume, err := bs.CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer bs.DeleteVolume(t, blockClient, volume)

	client.Microversion = "2.79"
	volumeAttachment, err := CreateVolumeAttachment(t, client, blockClient, server, volume)
	th.AssertNoErr(t, err)
	defer DeleteVolumeAttachment(t, client, blockClient, server, volumeAttachment)

	tools.PrintResource(t, volumeAttachment)

	th.AssertEquals(t, volumeAttachment.ServerID, server.ID)
}
