//go:build acceptance || blockstorage || volumeattachments

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	compute "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/compute/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumeAttachments(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	// minimu required microversion for volume attachments is 3.27
	blockClient.Microversion = "3.27"

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := compute.CreateServer(t, computeClient)
	th.AssertNoErr(t, err)
	defer compute.DeleteServer(t, computeClient, server)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = CreateVolumeAttachment(t, blockClient, volume, server)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(context.TODO(), blockClient, volume.ID).Extract()
	th.AssertNoErr(t, err)

	DeleteVolumeAttachment(t, blockClient, newVolume)
}
