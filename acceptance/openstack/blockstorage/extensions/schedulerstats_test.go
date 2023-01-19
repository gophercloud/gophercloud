//go:build acceptance || blockstorage
// +build acceptance blockstorage

package extensions

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	th "github.com/bizflycloud/gophercloud/testhelper"
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
