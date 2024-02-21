// Package v2 contains common functions for creating block storage based
// resources for use in acceptance tests. See the `*_test.go` files for
// example usages.
package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/backups"
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

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = snapshots.WaitForStatus(ctx, client, snapshot.ID, "available")
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

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, volume.ID, "available")
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

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, volume.ID, "available")
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
	err = tools.WaitFor(func(ctx context.Context) (bool, error) {
		_, err := snapshots.Get(ctx, client, snapshot.ID).Extract()
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

// CreateBackup will create a backup based on a volume. An error will be
// will be returned if the backup could not be created.
func CreateBackup(t *testing.T, client *gophercloud.ServiceClient, volumeID string) (*backups.Backup, error) {
	t.Logf("Attempting to create a backup of volume %s", volumeID)

	backupName := tools.RandomString("ACPTTEST", 16)
	createOpts := backups.CreateOpts{
		VolumeID: volumeID,
		Name:     backupName,
	}

	backup, err := backups.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	err = WaitForBackupStatus(client, backup.ID, "available")
	if err != nil {
		return nil, err
	}

	backup, err = backups.Get(context.TODO(), client, backup.ID).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created backup %s", backup.ID)
	tools.PrintResource(t, backup)

	th.AssertEquals(t, backup.Name, backupName)

	return backup, nil
}

// DeleteBackup will delete a backup. A fatal error will occur if the backup
// could not be deleted. This works best when used as a deferred function.
func DeleteBackup(t *testing.T, client *gophercloud.ServiceClient, backupID string) {
	if err := backups.Delete(context.TODO(), client, backupID).ExtractErr(); err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			t.Logf("Backup %s is already deleted", backupID)
			return
		}
		t.Fatalf("Unable to delete backup %s: %s", backupID, err)
	}

	t.Logf("Deleted backup %s", backupID)
}

// WaitForBackupStatus will continually poll a backup, checking for a particular
// status. It will do this for the amount of seconds defined.
func WaitForBackupStatus(client *gophercloud.ServiceClient, id, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		current, err := backups.Get(ctx, client, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok && status == "deleted" {
				return true, nil
			}
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}

// ResetBackupStatus will reset the status of a backup.
func ResetBackupStatus(t *testing.T, client *gophercloud.ServiceClient, backup *backups.Backup, status string) error {
	t.Logf("Attempting to reset the status of backup %s from %s to %s", backup.ID, backup.Status, status)

	resetOpts := backups.ResetStatusOpts{
		Status: status,
	}
	err := backups.ResetStatus(context.TODO(), client, backup.ID, resetOpts).ExtractErr()
	if err != nil {
		return err
	}

	return WaitForBackupStatus(client, backup.ID, status)
}
