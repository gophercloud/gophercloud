// Package v2 contains common functions for creating block storage based
// resources for use in acceptance tests. See the `*_test.go` files for
// example usages.
package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateSnapshot will create a snapshot of the specified volume.
// Snapshot will be assigned a random name and description.
func CreateSnapshot(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) (*snapshots.Snapshot, error) {
	snapshotName := tools.RandomString("ACPTTEST", 16)
	snapshotDescription := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create snapshot: %s", snapshotName)

	createOpts := snapshots.CreateOpts{
		VolumeID:    volume.ID,
		Name:        snapshotName,
		Description: snapshotDescription,
	}

	snapshot, err := snapshots.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return snapshot, err
	}

	err = snapshots.WaitForStatus(context.TODO(), client, snapshot.ID, "available", 60)
	if err != nil {
		return snapshot, err
	}

	t.Logf("Successfully created snapshot: %s", snapshot.ID)

	return snapshot, nil
}

// CreateVolume will create a volume with a random name and size of 1GB. An
// error will be returned if the volume was unable to be created.
func CreateVolume(t *testing.T, client *gophercloud.ServiceClient) (*volumes.Volume, error) {
	volumeName := tools.RandomString("ACPTTEST", 16)
	volumeDescription := tools.RandomString("ACPTTEST-DESC", 16)
	t.Logf("Attempting to create volume: %s", volumeName)

	createOpts := volumes.CreateOpts{
		Size:        1,
		Name:        volumeName,
		Description: volumeDescription,
	}

	volume, err := volumes.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return volume, err
	}

	err = volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60)
	if err != nil {
		return volume, err
	}

	tools.PrintResource(t, volume)
	th.AssertEquals(t, volume.Name, volumeName)
	th.AssertEquals(t, volume.Description, volumeDescription)
	th.AssertEquals(t, volume.Size, 1)

	t.Logf("Successfully created volume: %s", volume.ID)

	return volume, nil
}

// CreateVolumeFromImage will create a volume from with a random name and size of
// 1GB. An error will be returned if the volume was unable to be created.
func CreateVolumeFromImage(t *testing.T, client *gophercloud.ServiceClient) (*volumes.Volume, error) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	volumeName := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create volume: %s", volumeName)

	createOpts := volumes.CreateOpts{
		Size:    1,
		Name:    volumeName,
		ImageID: choices.ImageID,
	}

	volume, err := volumes.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return volume, err
	}

	err = volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60)
	if err != nil {
		return volume, err
	}

	newVolume, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return nil, err
	}

	th.AssertEquals(t, newVolume.Name, volumeName)
	th.AssertEquals(t, newVolume.Size, 1)

	t.Logf("Successfully created volume from image: %s", newVolume.ID)

	return newVolume, nil
}

// DeleteVolume will delete a volume. A fatal error will occur if the volume
// failed to be deleted. This works best when used as a deferred function.
func DeleteVolume(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) {
	t.Logf("Attempting to delete volume: %s", volume.ID)

	err := volumes.Delete(context.TODO(), client, volume.ID, volumes.DeleteOpts{}).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete volume %s: %v", volume.ID, err)
	}

	t.Logf("Successfully deleted volume: %s", volume.ID)
}

// DeleteSnapshot will delete a snapshot. A fatal error will occur if the
// snapshot failed to be deleted.
func DeleteSnapshot(t *testing.T, client *gophercloud.ServiceClient, snapshot *snapshots.Snapshot) {
	t.Logf("Attempting to delete snapshot: %s", snapshot.ID)

	err := snapshots.Delete(context.TODO(), client, snapshot.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete snapshot %s: %+v", snapshot.ID, err)
	}

	// Volumes can't be deleted until their snapshots have been,
	// so block until the snapshot is deleted.
	err = tools.WaitFor(func() (bool, error) {
		_, err := snapshots.Get(context.TODO(), client, snapshot.ID).Extract()
		if err != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		t.Fatalf("Error waiting for snapshot to delete: %v", err)
	}

	t.Logf("Successfully deleted snapshot: %s", snapshot.ID)
}
