//go:build acceptance || blockstorage
// +build acceptance blockstorage

package extensions

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestSchedulerStatsList(t *testing.T) {
	clients.RequireAdmin(t)

	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	listOpts := schedulerstats.ListOpts{
		Detail: true,
	}

	allPages, err := schedulerstats.List(blockClient, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allStats, err := schedulerstats.ExtractStoragePools(allPages)
	th.AssertNoErr(t, err)

	for _, stat := range allStats {
		tools.PrintResource(t, stat)
	}
}
