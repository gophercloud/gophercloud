//go:build acceptance || compute || volumeattach

package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumeAttachmentsListSinglePage(t *testing.T) {
	clients.RequireLong(t)
	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Register volume cleanup before server/attachment cleanup so dependent
	// resources are removed first, including when a later attachment fails.
	var expected []string
	for range 2 {
		volume, err := volumes.Create(ctx, blockClient, volumes.CreateOpts{
			Name: tools.RandomString("ACPTTEST-single-page-", 16),
			Size: 1,
		}, nil).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, volumes.Delete(ctx, blockClient, volume.ID, nil).ExtractErr())
		})
		th.AssertNoErr(t, volumes.WaitForStatus(ctx, blockClient, volume.ID, "available"))
		expected = append(expected, volume.ID)
	}
	server, _ := createPaginationServer(t, ctx, computeClient)
	for _, volumeID := range expected {
		_, err := volumeattach.Create(ctx, computeClient, server.ID, volumeattach.CreateOpts{VolumeID: volumeID}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			err := volumeattach.Delete(ctx, computeClient, server.ID, volumeID).ExtractErr()
			if err != nil {
				t.Errorf("Unable to detach volume %s: %v", volumeID, err)
				return
			}
			th.CheckNoErr(t, volumes.WaitForStatus(ctx, blockClient, volumeID, "available"))
		})
		th.AssertNoErr(t, volumes.WaitForStatus(ctx, blockClient, volumeID, "in-use"))
	}

	// Do not impose limit=1: this API intentionally returns one unpaginated
	// response. Both created volumes must be present in that response.
	tools.CheckSinglePage(t, ctx, computeClient, computeClient.ServiceURL("servers", server.ID, "os-volume_attachments"), func() pagination.Pager {
		return volumeattach.List(computeClient, server.ID)
	}, func(page pagination.Page) ([]string, error) {
		items, err := volumeattach.ExtractVolumeAttachments(page)
		ids := make([]string, len(items))
		for i, item := range items {
			th.CheckEquals(t, server.ID, item.ServerID)
			ids[i] = item.VolumeID
		}
		return ids, err
	}, expected)
}
