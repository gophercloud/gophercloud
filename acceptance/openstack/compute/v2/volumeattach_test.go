// +build acceptance compute volumeattach

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	bs "github.com/gophercloud/gophercloud/acceptance/openstack/blockstorage/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestVolumeAttachAttachment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	blockClient, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	server, err := CreateServer(t, client)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer DeleteServer(t, client, server)

	volume, err := bs.CreateVolume(t, blockClient)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer bs.DeleteVolume(t, blockClient, volume)

	volumeAttachment, err := CreateVolumeAttachment(t, client, blockClient, server, volume)
	if err != nil {
		t.Fatalf("Unable to attach volume: %v", err)
	}
	defer DeleteVolumeAttachment(t, client, blockClient, server, volumeAttachment)

	tools.PrintResource(t, volumeAttachment)

	th.AssertEquals(t, volumeAttachment.ServerID, server.ID)
}
