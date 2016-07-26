// +build acceptance compute volumeattach

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestVolumeAttachAttachment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	blockClient, err := newBlockClient()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	volume, err := createVolume(t, blockClient)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}

	if err = volumes.WaitForStatus(blockClient, volume.ID, "available", 60); err != nil {
		t.Fatalf("Unable to wait for volume: %v", err)
	}
	defer deleteVolume(t, blockClient, volume)

	volumeAttachment, err := createVolumeAttachment(t, client, blockClient, server, volume)
	if err != nil {
		t.Fatalf("Unable to attach volume: %v", err)
	}
	defer deleteVolumeAttachment(t, client, blockClient, server, volumeAttachment)

	printVolumeAttachment(t, volumeAttachment)

}

func createVolume(t *testing.T, blockClient *gophercloud.ServiceClient) (*volumes.Volume, error) {
	volumeName := tools.RandomString("ACPTTEST", 16)
	createOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume, err := volumes.Create(blockClient, createOpts).Extract()
	if err != nil {
		return volume, err
	}

	t.Logf("Created volume: %s", volume.ID)
	return volume, nil
}

func deleteVolume(t *testing.T, blockClient *gophercloud.ServiceClient, volume *volumes.Volume) {
	err := volumes.Delete(blockClient, volume.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete volume: %v", err)
	}

	t.Logf("Deleted volume: %s", volume.ID)
}

func createVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, blockClient *gophercloud.ServiceClient, server *servers.Server, volume *volumes.Volume) (*volumeattach.VolumeAttachment, error) {
	volumeAttachOptions := volumeattach.CreateOpts{
		VolumeID: volume.ID,
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)
	volumeAttachment, err := volumeattach.Create(client, server.ID, volumeAttachOptions).Extract()
	if err != nil {
		return volumeAttachment, err
	}

	if err = volumes.WaitForStatus(blockClient, volume.ID, "in-use", 60); err != nil {
		return volumeAttachment, err
	}

	return volumeAttachment, nil
}

func deleteVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, blockClient *gophercloud.ServiceClient, server *servers.Server, volumeAttachment *volumeattach.VolumeAttachment) {

	err := volumeattach.Delete(client, server.ID, volumeAttachment.VolumeID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to detach volume: %v", err)
	}

	if err = volumes.WaitForStatus(blockClient, volumeAttachment.ID, "available", 60); err != nil {
		t.Fatalf("Unable to wait for volume: %v", err)
	}
	t.Logf("Deleted volume: %s", volumeAttachment.VolumeID)
}

func printVolumeAttachment(t *testing.T, volumeAttachment *volumeattach.VolumeAttachment) {
	t.Logf("ID: %s", volumeAttachment.ID)
	t.Logf("Device: %s", volumeAttachment.Device)
	t.Logf("VolumeID: %s", volumeAttachment.VolumeID)
	t.Logf("ServerID: %s", volumeAttachment.ServerID)
}
