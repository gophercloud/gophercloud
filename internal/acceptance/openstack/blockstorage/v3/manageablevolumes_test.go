//go:build acceptance || blockstorage || volumes

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestManageableVolumes(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	client.Microversion = "3.8"

	volume1, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)

	err = Unmanage(t, client, volume1)
	if err != nil {
		DeleteVolume(t, client, volume1)
	}
	th.AssertNoErr(t, err)

	managed1, err := ManageExisting(t, client, volume1)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, managed1)

	th.CheckEquals(t, volume1.Host, managed1.Host)
	th.AssertEquals(t, volume1.Name, managed1.Name)
	th.AssertEquals(t, volume1.AvailabilityZone, managed1.AvailabilityZone)
	th.AssertEquals(t, volume1.Description, managed1.Description)
	th.AssertEquals(t, volume1.VolumeType, managed1.VolumeType)
	th.AssertEquals(t, volume1.Bootable, managed1.Bootable)
	th.AssertDeepEquals(t, volume1.Metadata, managed1.Metadata)
	th.AssertEquals(t, volume1.Size, managed1.Size)

	allPages, err := volumes.List(client, volumes.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allVolumes, err := volumes.ExtractVolumes(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allVolumes {
		if v.ID == managed1.ID {
			found = true
			break
		}
	}
	th.AssertEquals(t, true, found)
}
