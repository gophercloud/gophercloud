package v3

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/attachments"
	v3 "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// CreateVolumeAttachment will attach a volume to an instance. An error will be
// returned if the attachment failed.
func CreateVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume, server *servers.Server) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume attachment in short mode.")
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)

	var err error
	var attachment *volumeattach.VolumeAttachment
	if attachment, err = volumeattach.Create(client, server.ID, volumeattach.CreateOpts{
		VolumeID: volume.ID,
	}).Extract(); err != nil {
		return err
	}

	// The nova API doesn't provide a way to check whether the volume is really attached so
	// we use the Cinder API to wait for the volume to become attached.
	mv := client.Microversion
	client.Microversion = "3.44"
	defer func() {
		client.Microversion = mv
	}()
	if err = attachments.Complete(client, attachment.ID).ExtractErr(); err != nil {
		return err
	}
	if err = attachments.WaitForStatus(client, attachment.ID, "attached", 60); err != nil {
		e := attachments.Delete(client, attachment.ID).ExtractErr()
		if e != nil {
			t.Logf("Failed to delete %q attachment: %s", attachment.ID, err)
		}
		return err
	}

	attachment, err = volumeattach.Get(client, server.ID, attachment.ID).Extract()
	if err != nil {
		return err
	}

	allPages, err := volumeattach.List(client, server.ID).AllPages()
	if err != nil {
		return err
	}

	AllAttachments, err := volumeattach.ExtractVolumeAttachments(allPages)
	if err != nil {
		return err
	}

	if AllAttachments[0].ID != attachment.ID {
		return fmt.Errorf("attachment IDs from get and list are not equal: %q != %q", AllAttachments[0].ID, attachment.ID)
	}

	t.Logf("Attached volume %s to server %s within %q attachment", volume.ID, server.ID, attachment.ID)

	return nil
}

// DeleteVolumeAttachment will detach a volume from an instance. A fatal error
// will occur if the attachment failed to be deleted.
func DeleteVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume) {
	t.Logf("Attepting to detach volume volume: %s", volume.ID)

	if err := volumeattach.Delete(client, volume.Attachments[0].ServerID, volume.Attachments[0].AttachmentID).ExtractErr(); err != nil {
		t.Fatalf("Unable to detach volume %s: %v", volume.ID, err)
	}

	if err := v3.WaitForStatus(client, volume.ID, "available", 60); err != nil {
		t.Fatalf("Volume %s failed to become unavailable in 60 seconds: %v", volume.ID, err)
	}

	t.Logf("Detached volume: %s", volume.ID)
}
