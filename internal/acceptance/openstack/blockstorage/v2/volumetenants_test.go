//go:build acceptance || blockstorage || volumetenants

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVolumeTenants(t *testing.T) {
	var allVolumes []volumes.Volume

	clients.SkipReleasesAbove(t, "stable/ocata")

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	listOpts := volumes.ListOpts{
		Name: "I SHOULD NOT EXIST",
	}
	allPages, err := volumes.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = volumes.ExtractVolumesInto(allPages, &allVolumes)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allVolumes))

	volume1, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume1)

	allPages, err = volumes.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = volumes.ExtractVolumesInto(allPages, &allVolumes)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, len(allVolumes) > 0)
}
