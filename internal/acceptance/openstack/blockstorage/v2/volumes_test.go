//go:build acceptance || blockstorage || volumes

package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	compute "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/compute/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumesCreateDestroy(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	newVolume, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	th.AssertNoErr(t, err)

	// Update volume
	updatedVolumeName := ""
	updatedVolumeDescription := ""
	updateOpts := volumes.UpdateOpts{
		Name:        &updatedVolumeName,
		Description: &updatedVolumeDescription,
	}
	updatedVolume, err := volumes.Update(context.TODO(), client, volume.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedVolume)
	th.AssertEquals(t, updatedVolume.Name, updatedVolumeName)
	th.AssertEquals(t, updatedVolume.Description, updatedVolumeDescription)

	allPages, err := volumes.List(client, volumes.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allVolumes, err := volumes.ExtractVolumes(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allVolumes {
		tools.PrintResource(t, volume)
		if v.ID == newVolume.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestVolumesCreateForceDestroy(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newVolume)

	err = volumes.ForceDelete(context.TODO(), client, newVolume.ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVolumesCascadeDelete(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	vol, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, vol.ID, "available")
	th.AssertNoErr(t, err)

	snapshot1, err := CreateSnapshot(t, client, vol)
	th.AssertNoErr(t, err)

	snapshot2, err := CreateSnapshot(t, client, vol)
	th.AssertNoErr(t, err)

	t.Logf("Attempting to delete volume: %s", vol.ID)

	deleteOpts := volumes.DeleteOpts{Cascade: true}
	err = volumes.Delete(context.TODO(), client, vol.ID, deleteOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete volume %s: %v", vol.ID, err)
	}

	for _, sid := range []string{snapshot1.ID, snapshot2.ID} {
		err := tools.WaitFor(func(ctx context.Context) (bool, error) {
			_, err := snapshots.Get(ctx, client, sid).Extract()
			if err != nil {
				return true, nil
			}
			return false, nil
		})
		th.AssertNoErr(t, err)
		t.Logf("Successfully deleted snapshot: %s", sid)
	}

	err = tools.WaitFor(func(ctx context.Context) (bool, error) {
		_, err := volumes.Get(ctx, client, vol.ID).Extract()
		if err != nil {
			return true, nil
		}
		return false, nil
	})
	th.AssertNoErr(t, err)

	t.Logf("Successfully deleted volume: %s", vol.ID)
}

func TestVolumeActionsUploadImageDestroy(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	imageClient, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	volumeImage, err := CreateUploadImage(t, blockClient, volume)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, volumeImage)

	err = DeleteUploadedImage(t, imageClient, volumeImage.ImageID)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsAttachCreateDestroy(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := compute.CreateServer(t, computeClient)
	th.AssertNoErr(t, err)
	defer compute.DeleteServer(t, computeClient, server)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = CreateVolumeAttach(t, blockClient, volume, server)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(context.TODO(), blockClient, volume.ID).Extract()
	th.AssertNoErr(t, err)

	DeleteVolumeAttach(t, blockClient, newVolume)
}

func TestVolumeActionsReserveUnreserve(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	err = CreateVolumeReserve(t, client, volume)
	th.AssertNoErr(t, err)
	defer DeleteVolumeReserve(t, client, volume)
}

func TestVolumeActionsExtendSize(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	tools.PrintResource(t, volume)

	err = ExtendVolumeSize(t, blockClient, volume)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(context.TODO(), blockClient, volume.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newVolume)
}

func TestVolumeActionsImageMetadata(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = SetImageMetadata(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsSetBootable(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = SetBootable(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsResetStatus(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	tools.PrintResource(t, volume)

	err = ResetVolumeStatus(t, client, volume, "error")
	th.AssertNoErr(t, err)

	err = ResetVolumeStatus(t, client, volume, "available")
	th.AssertNoErr(t, err)
}

// Note(jtopjian): I plan to work on this at some point, but it requires
// setting up a server with iscsi utils.
/*
func TestVolumeConns(t *testing.T) {
    clients.SkipReleasesAbove(t, "stable/ocata")

    client, err := newClient()
    th.AssertNoErr(t, err)

    t.Logf("Creating volume")
    cv, err := volumes.Create(client, &volumes.CreateOpts{
        Size: 1,
        Name: "blockv2-volume",
    }, nil).Extract()
    th.AssertNoErr(t, err)

    defer func() {
        err = volumes.WaitForStatus(client, cv.ID, "available", 60)
        th.AssertNoErr(t, err)

        t.Logf("Deleting volume")
        err = volumes.Delete(client, cv.ID, volumes.DeleteOpts{}).ExtractErr()
        th.AssertNoErr(t, err)
    }()

    err = volumes.WaitForStatus(client, cv.ID, "available", 60)
    th.AssertNoErr(t, err)

    connOpts := &volumes.ConnectorOpts{
        IP:        "127.0.0.1",
        Host:      "stack",
        Initiator: "iqn.1994-05.com.redhat:17cf566367d2",
        Multipath: false,
        Platform:  "x86_64",
        OSType:    "linux2",
    }

    t.Logf("Initializing connection")
    _, err = volumes.InitializeConnection(client, cv.ID, connOpts).Extract()
    th.AssertNoErr(t, err)

    t.Logf("Terminating connection")
    err = volumes.TerminateConnection(client, cv.ID, connOpts).ExtractErr()
    th.AssertNoErr(t, err)
}
*/
