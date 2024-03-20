//go:build acceptance || blockstorage || snapshots

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/snapshots"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSnapshots(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume)

	snapshot, err := CreateSnapshot(t, client, volume)
	th.AssertNoErr(t, err)
	defer DeleteSnapshot(t, client, snapshot)

	newSnapshot, err := snapshots.Get(context.TODO(), client, snapshot.ID).Extract()
	th.AssertNoErr(t, err)

	allPages, err := snapshots.List(client, snapshots.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allSnapshots {
		tools.PrintResource(t, snapshot)
		if v.ID == newSnapshot.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
