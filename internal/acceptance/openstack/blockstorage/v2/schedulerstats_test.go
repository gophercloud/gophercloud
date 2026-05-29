//go:build acceptance || blockstorage || schedulerstats

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/schedulerstats"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSchedulerStatsList(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireAdmin(t)

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	listOpts := schedulerstats.ListOpts{
		Detail: true,
	}

	allPages, err := schedulerstats.List(blockClient, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allStats, err := schedulerstats.ExtractStoragePools(allPages)
	th.AssertNoErr(t, err)

	for _, stat := range allStats {
		tools.PrintResource(t, stat)
	}
}
