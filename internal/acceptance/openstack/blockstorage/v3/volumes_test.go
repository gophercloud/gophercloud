//go:build acceptance || blockstorage
// +build acceptance blockstorage

package v3

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
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

	vol, err := volumes.Create(context.TODO(), client, volOpts).Extract()
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
