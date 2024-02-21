//go:build acceptance || blockstorage || volumes

package v3

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	compute "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/compute/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumes(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume1, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume1)

	volume2, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume2)

	// Update volume
	updatedVolumeName := ""
	updatedVolumeDescription := ""
	updateOpts := volumes.UpdateOpts{
		Name:        &updatedVolumeName,
		Description: &updatedVolumeDescription,
	}
	updatedVolume, err := volumes.Update(context.TODO(), client, volume1.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedVolume)
	th.AssertEquals(t, updatedVolume.Name, updatedVolumeName)
	th.AssertEquals(t, updatedVolume.Description, updatedVolumeDescription)

	listOpts := volumes.ListOpts{
		Limit: 1,
	}

	err = volumes.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		actual, err := volumes.ExtractVolumes(page)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(actual))

		var found bool
		for _, v := range actual {
			if v.ID == volume1.ID || v.ID == volume2.ID {
				found = true
			}
		}

		th.AssertEquals(t, found, true)

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func TestVolumesMultiAttach(t *testing.T) {
	clients.RequireAdmin(t)
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	vt, err := CreateVolumeTypeMultiAttach(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	volumeName := tools.RandomString("ACPTTEST", 16)

	volOpts := volumes.CreateOpts{
		Size:        1,
		Name:        volumeName,
		Description: "Testing creation of multiattach enabled volume",
		VolumeType:  vt.ID,
	}

	vol, err := volumes.Create(context.TODO(), client, volOpts, nil).Extract()
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, vol)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, vol.ID, "available")
	th.AssertNoErr(t, err)

	th.AssertEquals(t, vol.Multiattach, true)
}

func TestVolumesCascadeDelete(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
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
	blockClient, err := clients.NewBlockStorageV3Client()
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
	blockClient, err := clients.NewBlockStorageV3Client()
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
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	err = CreateVolumeReserve(t, client, volume)
	th.AssertNoErr(t, err)
	defer DeleteVolumeReserve(t, client, volume)
}

func TestVolumeActionsExtendSize(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
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
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = SetImageMetadata(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsSetBootable(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = SetBootable(t, blockClient, volume)
	th.AssertNoErr(t, err)
}

func TestVolumeActionsChangeType(t *testing.T) {
	//	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volumeType1, err := CreateVolumeTypeNoExtraSpecs(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, volumeType1)

	volumeType2, err := CreateVolumeTypeNoExtraSpecs(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, volumeType2)

	volume, err := CreateVolumeWithType(t, client, volumeType1)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	tools.PrintResource(t, volume)

	err = ChangeVolumeType(t, client, volume, volumeType2)
	th.AssertNoErr(t, err)

	newVolume, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newVolume.VolumeType, volumeType2.Name)

	tools.PrintResource(t, newVolume)
}

func TestVolumeActionsResetStatus(t *testing.T) {
	client, err := clients.NewBlockStorageV3Client()
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

func TestVolumeActionsReImage(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/yoga")

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)
	blockClient.Microversion = "3.68"

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	err = ReImage(t, blockClient, volume, choices.ImageID)
	th.AssertNoErr(t, err)
}

// Note(jtopjian): I plan to work on this at some point, but it requires
// setting up a server with iscsi utils.
/*
func TestVolumeConns(t *testing.T) {
    client, err := newClient()
    th.AssertNoErr(t, err)

    t.Logf("Creating volume")
    cv, err := volumes.Create(context.TODO(), client, &volumes.CreateOpts{
        Size: 1,
        Name: "blockv2-volume",
    }, nil).Extract()
    th.AssertNoErr(t, err)

    defer func() {
        err = volumes.WaitForStatus(context.TODO(), client, cv.ID, "available", 60)
        th.AssertNoErr(t, err)

        t.Logf("Deleting volume")
        err = volumes.Delete(context.TODO(), client, cv.ID, volumes.DeleteOpts{}).ExtractErr()
        th.AssertNoErr(t, err)
    }()

    err = volumes.WaitForStatus(context.TODO(), client, cv.ID, "available", 60)
    th.AssertNoErr(t, err)

    connOpts := &ConnectorOpts{
        IP:        "127.0.0.1",
        Host:      "stack",
        Initiator: "iqn.1994-05.com.redhat:17cf566367d2",
        Multipath: false,
        Platform:  "x86_64",
        OSType:    "linux2",
    }

    t.Logf("Initializing connection")
    _, err = InitializeConnection(client, cv.ID, connOpts).Extract()
    th.AssertNoErr(t, err)

    t.Logf("Terminating connection")
    err = TerminateConnection(client, cv.ID, connOpts).ExtractErr()
    th.AssertNoErr(t, err)
}
*/
