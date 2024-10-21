// Package v2 contains common functions for creating block storage based
// resources for use in acceptance tests. See the `*_test.go` files for
// example usages.
package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/backups"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
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

	volume, err := volumes.Create(context.TODO(), client, createOpts, nil).Extract()
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

	volume, err := volumes.Create(context.TODO(), client, createOpts, nil).Extract()
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
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) && status == "deleted" {
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

// CreateUploadImage will upload volume it as volume-baked image. An name of new image or err will be
// returned
func CreateUploadImage(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) (volumes.VolumeImage, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires volume-backed image uploading in short mode.")
	}

	imageName := tools.RandomString("ACPTTEST", 16)
	uploadImageOpts := volumes.UploadImageOpts{
		ImageName: imageName,
		Force:     true,
	}

	volumeImage, err := volumes.UploadImage(context.TODO(), client, volume.ID, uploadImageOpts).Extract()
	if err != nil {
		return volumeImage, err
	}

	t.Logf("Uploading volume %s as volume-backed image %s", volume.ID, imageName)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := volumes.WaitForStatus(ctx, client, volume.ID, "available"); err != nil {
		return volumeImage, err
	}

	t.Logf("Uploaded volume %s as volume-backed image %s", volume.ID, imageName)

	return volumeImage, nil

}

// DeleteUploadedImage deletes uploaded image. An error will be returned
// if the deletion request failed.
func DeleteUploadedImage(t *testing.T, client *gophercloud.ServiceClient, imageID string) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume-backed image removing in short mode.")
	}

	t.Logf("Removing image %s", imageID)

	err := images.Delete(context.TODO(), client, imageID).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

// CreateVolumeAttach will attach a volume to an instance. An error will be
// returned if the attachment failed.
func CreateVolumeAttach(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume, server *servers.Server) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume attachment in short mode.")
	}

	attachOpts := volumes.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: server.ID,
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)

	if err := volumes.Attach(context.TODO(), client, volume.ID, attachOpts).ExtractErr(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := volumes.WaitForStatus(ctx, client, volume.ID, "in-use"); err != nil {
		return err
	}

	t.Logf("Attached volume %s to server %s", volume.ID, server.ID)

	return nil
}

// CreateVolumeReserve creates a volume reservation. An error will be returned
// if the reservation failed.
func CreateVolumeReserve(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume reservation in short mode.")
	}

	t.Logf("Attempting to reserve volume %s", volume.ID)

	if err := volumes.Reserve(context.TODO(), client, volume.ID).ExtractErr(); err != nil {
		return err
	}

	t.Logf("Reserved volume %s", volume.ID)

	return nil
}

// DeleteVolumeAttach will detach a volume from an instance. A fatal error will
// occur if the snapshot failed to be deleted. This works best when used as a
// deferred function.
func DeleteVolumeAttach(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) {
	t.Logf("Attepting to detach volume volume: %s", volume.ID)

	detachOpts := volumes.DetachOpts{
		AttachmentID: volume.Attachments[0].AttachmentID,
	}

	if err := volumes.Detach(context.TODO(), client, volume.ID, detachOpts).ExtractErr(); err != nil {
		t.Fatalf("Unable to detach volume %s: %v", volume.ID, err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := volumes.WaitForStatus(ctx, client, volume.ID, "available"); err != nil {
		t.Fatalf("Volume %s failed to become unavailable in 60 seconds: %v", volume.ID, err)
	}

	t.Logf("Detached volume: %s", volume.ID)
}

// DeleteVolumeReserve deletes a volume reservation. A fatal error will occur
// if the deletion request failed. This works best when used as a deferred
// function.
func DeleteVolumeReserve(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) {
	if testing.Short() {
		t.Skip("Skipping test that requires volume reservation in short mode.")
	}

	t.Logf("Attempting to unreserve volume %s", volume.ID)

	if err := volumes.Unreserve(context.TODO(), client, volume.ID).ExtractErr(); err != nil {
		t.Fatalf("Unable to unreserve volume %s: %v", volume.ID, err)
	}

	t.Logf("Unreserved volume %s", volume.ID)
}

// ExtendVolumeSize will extend the size of a volume.
func ExtendVolumeSize(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to extend the size of volume %s", volume.ID)

	extendOpts := volumes.ExtendSizeOpts{
		NewSize: 2,
	}

	err := volumes.ExtendSize(context.TODO(), client, volume.ID, extendOpts).ExtractErr()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := volumes.WaitForStatus(ctx, client, volume.ID, "available"); err != nil {
		return err
	}

	return nil
}

// SetImageMetadata will apply the metadata to a volume.
func SetImageMetadata(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply image metadata to volume %s", volume.ID)

	imageMetadataOpts := volumes.ImageMetadataOpts{
		Metadata: map[string]string{
			"image_name": "testimage",
		},
	}

	err := volumes.SetImageMetadata(context.TODO(), client, volume.ID, imageMetadataOpts).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

// SetBootable will set a bootable status to a volume.
func SetBootable(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply bootable status to volume %s", volume.ID)

	bootableOpts := volumes.BootableOpts{
		Bootable: true,
	}

	err := volumes.SetBootable(context.TODO(), client, volume.ID, bootableOpts).ExtractErr()
	if err != nil {
		return err
	}

	vol, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) != "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'true'", vol.Bootable)
	}

	bootableOpts = volumes.BootableOpts{
		Bootable: false,
	}

	err = volumes.SetBootable(context.TODO(), client, volume.ID, bootableOpts).ExtractErr()
	if err != nil {
		return err
	}

	vol, err = volumes.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) == "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'false'", vol.Bootable)
	}

	return nil
}

// ResetVolumeStatus will reset the status of a volume.
func ResetVolumeStatus(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume, status string) error {
	t.Logf("Attempting to reset the status of volume %s from %s to %s", volume.ID, volume.Status, status)

	resetOpts := volumes.ResetStatusOpts{
		Status: status,
	}
	err := volumes.ResetStatus(context.TODO(), client, volume.ID, resetOpts).ExtractErr()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := volumes.WaitForStatus(ctx, client, volume.ID, status); err != nil {
		return err
	}

	return nil
}
