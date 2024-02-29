package v3

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/attachments"
	v3 "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)

// CreateVolumeAttachment will attach a volume to an instance. An error will be
// returned if the attachment failed.
func CreateVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume, server *servers.Server) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume attachment in short mode.")
	}

	attachOpts := &attachments.CreateOpts{
		VolumeUUID:   volume.ID,
		InstanceUUID: server.ID,
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)

	var err error
	var attachment *attachments.Attachment
	if attachment, err = attachments.Create(context.TODO(), client, attachOpts).Extract(); err != nil {
		return err
	}

	mv := client.Microversion
	client.Microversion = "3.44"
	defer func() {
		client.Microversion = mv
	}()
	if err = attachments.Complete(context.TODO(), client, attachment.ID).ExtractErr(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err = attachments.WaitForStatus(ctx, client, attachment.ID, "attached"); err != nil {
		e := attachments.Delete(context.TODO(), client, attachment.ID).ExtractErr()
		if e != nil {
			t.Logf("Failed to delete %q attachment: %s", attachment.ID, err)
		}
		return err
	}

	attachment, err = attachments.Get(context.TODO(), client, attachment.ID).Extract()
	if err != nil {
		return err
	}

	listOpts := &attachments.ListOpts{
		VolumeID:   volume.ID,
		InstanceID: server.ID,
	}
	allPages, err := attachments.List(client, listOpts).AllPages(context.TODO())
	if err != nil {
		return err
	}

	allAttachments, err := attachments.ExtractAttachments(allPages)
	if err != nil {
		return err
	}

	if allAttachments[0].ID != attachment.ID {
		return fmt.Errorf("Attachment IDs from get and list are not equal: %q != %q", allAttachments[0].ID, attachment.ID)
	}

	t.Logf("Attached volume %s to server %s within %q attachment", volume.ID, server.ID, attachment.ID)

	return nil
}

// DeleteVolumeAttachment will detach a volume from an instance. A fatal error
// will occur if the attachment failed to be deleted.
func DeleteVolumeAttachment(t *testing.T, client *gophercloud.ServiceClient, volume *v3.Volume) {
	t.Logf("Attepting to detach volume volume: %s", volume.ID)

	if err := attachments.Delete(context.TODO(), client, volume.Attachments[0].AttachmentID).ExtractErr(); err != nil {
		t.Fatalf("Unable to detach volume %s: %v", volume.ID, err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if err := v3.WaitForStatus(ctx, client, volume.ID, "available"); err != nil {
		t.Fatalf("Volume %s failed to become unavailable in 60 seconds: %v", volume.ID, err)
	}

	t.Logf("Detached volume: %s", volume.ID)
}
