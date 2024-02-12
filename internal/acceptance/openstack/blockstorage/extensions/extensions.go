// Package extensions contains common functions for creating block storage
// resources that are extensions of the block storage API. See the `*_test.go`
// files for example usages.
package extensions

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/extensions/backups"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	v3 "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateUploadImage will upload volume it as volume-baked image. An name of new image or err will be
// returned
func CreateUploadImage(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) (volumeactions.VolumeImage, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires volume-backed image uploading in short mode.")
	}

	imageName := tools.RandomString("ACPTTEST", 16)
	uploadImageOpts := volumeactions.UploadImageOpts{
		ImageName: imageName,
		Force:     true,
	}

	volumeImage, err := volumeactions.UploadImage(context.TODO(), client, volume.ID, uploadImageOpts).Extract()
	if err != nil {
		return volumeImage, err
	}

	t.Logf("Uploading volume %s as volume-backed image %s", volume.ID, imageName)

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60); err != nil {
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

	attachOpts := volumeactions.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: server.ID,
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)

	if err := volumeactions.Attach(context.TODO(), client, volume.ID, attachOpts).ExtractErr(); err != nil {
		return err
	}

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, "in-use", 60); err != nil {
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

	if err := volumeactions.Reserve(context.TODO(), client, volume.ID).ExtractErr(); err != nil {
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

	detachOpts := volumeactions.DetachOpts{
		AttachmentID: volume.Attachments[0].AttachmentID,
	}

	if err := volumeactions.Detach(context.TODO(), client, volume.ID, detachOpts).ExtractErr(); err != nil {
		t.Fatalf("Unable to detach volume %s: %v", volume.ID, err)
	}

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60); err != nil {
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

	if err := volumeactions.Unreserve(context.TODO(), client, volume.ID).ExtractErr(); err != nil {
		t.Fatalf("Unable to unreserve volume %s: %v", volume.ID, err)
	}

	t.Logf("Unreserved volume %s", volume.ID)
}

// ExtendVolumeSize will extend the size of a volume.
func ExtendVolumeSize(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to extend the size of volume %s", volume.ID)

	extendOpts := volumeactions.ExtendSizeOpts{
		NewSize: 2,
	}

	err := volumeactions.ExtendSize(context.TODO(), client, volume.ID, extendOpts).ExtractErr()
	if err != nil {
		return err
	}

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60); err != nil {
		return err
	}

	return nil
}

// SetImageMetadata will apply the metadata to a volume.
func SetImageMetadata(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply image metadata to volume %s", volume.ID)

	imageMetadataOpts := volumeactions.ImageMetadataOpts{
		Metadata: map[string]string{
			"image_name": "testimage",
		},
	}

	err := volumeactions.SetImageMetadata(context.TODO(), client, volume.ID, imageMetadataOpts).ExtractErr()
	if err != nil {
		return err
	}

	return nil
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
	return tools.WaitFor(func() (bool, error) {
		current, err := backups.Get(context.TODO(), client, id).Extract()
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

// SetBootable will set a bootable status to a volume.
func SetBootable(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply bootable status to volume %s", volume.ID)

	bootableOpts := volumeactions.BootableOpts{
		Bootable: true,
	}

	err := volumeactions.SetBootable(context.TODO(), client, volume.ID, bootableOpts).ExtractErr()
	if err != nil {
		return err
	}

	vol, err := v3.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) != "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'true'", vol.Bootable)
	}

	bootableOpts = volumeactions.BootableOpts{
		Bootable: false,
	}

	err = volumeactions.SetBootable(context.TODO(), client, volume.ID, bootableOpts).ExtractErr()
	if err != nil {
		return err
	}

	vol, err = v3.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) == "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'false'", vol.Bootable)
	}

	return nil
}

// ChangeVolumeType will extend the size of a volume.
func ChangeVolumeType(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume, vt *volumetypes.VolumeType) error {
	t.Logf("Attempting to change the type of volume %s from %s to %s", volume.ID, volume.VolumeType, vt.Name)

	changeOpts := volumeactions.ChangeTypeOpts{
		NewType:         vt.Name,
		MigrationPolicy: volumeactions.MigrationPolicyOnDemand,
	}

	err := volumeactions.ChangeType(context.TODO(), client, volume.ID, changeOpts).ExtractErr()
	if err != nil {
		return err
	}

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60); err != nil {
		return err
	}

	return nil
}

// ResetVolumeStatus will reset the status of a volume.
func ResetVolumeStatus(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume, status string) error {
	t.Logf("Attempting to reset the status of volume %s from %s to %s", volume.ID, volume.Status, status)

	resetOpts := volumeactions.ResetStatusOpts{
		Status: status,
	}
	err := volumeactions.ResetStatus(context.TODO(), client, volume.ID, resetOpts).ExtractErr()
	if err != nil {
		return err
	}

	if err := volumes.WaitForStatus(context.TODO(), client, volume.ID, status, 60); err != nil {
		return err
	}

	return nil
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

// ReImage will re-image a volume
func ReImage(t *testing.T, client *gophercloud.ServiceClient, volume *volumes.Volume, imageID string) error {
	t.Logf("Attempting to re-image volume %s", volume.ID)

	reimageOpts := volumeactions.ReImageOpts{
		ImageID:         imageID,
		ReImageReserved: false,
	}

	err := volumeactions.ReImage(context.TODO(), client, volume.ID, reimageOpts).ExtractErr()
	if err != nil {
		return err
	}

	err = volumes.WaitForStatus(context.TODO(), client, volume.ID, "available", 60)
	if err != nil {
		return err
	}

	vol, err := v3.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		return err
	}

	if vol.VolumeImageMetadata == nil {
		return fmt.Errorf("volume does not have VolumeImageMetadata map")
	}

	if strings.ToLower(vol.VolumeImageMetadata["image_id"]) != imageID {
		return fmt.Errorf("volume image id '%s', expected '%s'", vol.VolumeImageMetadata["image_id"], imageID)
	}

	return nil
}
