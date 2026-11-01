//go:build acceptance || sharedfilesystems || schedulerstats

package v2

import (
	"context"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/schedulerstats"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSchedulerStatsList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.23"

	allPages, err := schedulerstats.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPools, err := schedulerstats.ExtractPools(allPages)
	th.AssertNoErr(t, err)

	for _, recordset := range allPools {
		tools.PrintResource(t, &recordset)
	}
}

func TestSchedulerStatsFilters(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.23"

	pages, err := schedulerstats.ListDetail(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	pools, err := schedulerstats.ExtractPools(pages)
	th.AssertNoErr(t, err)
	if len(pools) == 0 {
		t.Fatal("Need at least one storage pool to verify scheduler filters")
	}
	pool := pools[0]

	// A share type exercises capability filtering through Manila's supported API.
	shareType, err := sharetypes.Create(context.TODO(), client, sharetypes.CreateOpts{
		Name:       tools.RandomString("ACPTTEST", 16),
		ExtraSpecs: sharetypes.ExtraSpecsOpts{DriverHandlesShareServers: &pool.Capabilities.DriverHandlesShareServers},
	}).Extract()
	th.AssertNoErr(t, err)
	defer DeleteShareType(t, client, shareType)

	for _, detail := range []bool{false, true} {
		name := "summary"
		if detail {
			name = "detail"
		}
		t.Run(name, func(t *testing.T) {
			for _, test := range []struct {
				name string
				opts schedulerstats.ListOpts
			}{
				{"pool", schedulerstats.ListOpts{PoolName: "^" + regexp.QuoteMeta(pool.Pool) + "$"}},
				{"host", schedulerstats.ListOpts{HostName: "^" + regexp.QuoteMeta(pool.Host) + "$"}},
				{"backend", schedulerstats.ListOpts{BackendName: "^" + regexp.QuoteMeta(pool.Backend) + "$"}},
				{"share_type", schedulerstats.ListOpts{ShareType: shareType.ID}},
			} {
				t.Run(test.name, func(t *testing.T) {
					var pager pagination.Pager
					if detail {
						pager = schedulerstats.ListDetail(client, schedulerstats.ListDetailOpts(test.opts))
					} else {
						pager = schedulerstats.List(client, test.opts)
					}
					pages, err := pager.AllPages(context.TODO())
					th.AssertNoErr(t, err)
					filtered, err := schedulerstats.ExtractPools(pages)
					th.AssertNoErr(t, err)
					found := false
					for _, result := range filtered {
						if result.Name == pool.Name {
							found = true
						}
						switch test.name {
						case "pool":
							th.CheckEquals(t, pool.Pool, result.Pool)
						case "host":
							th.CheckEquals(t, pool.Host, result.Host)
						case "backend":
							th.CheckEquals(t, pool.Backend, result.Backend)
						}
					}
					th.AssertTrue(t, found)
				})
			}
		})
	}

	for name, pager := range map[string]pagination.Pager{
		"pool":           schedulerstats.List(client, schedulerstats.ListOpts{PoolName: "a^"}),
		"host":           schedulerstats.List(client, schedulerstats.ListOpts{HostName: "a^"}),
		"backend":        schedulerstats.List(client, schedulerstats.ListOpts{BackendName: "a^"}),
		"detail_pool":    schedulerstats.ListDetail(client, schedulerstats.ListDetailOpts{PoolName: "a^"}),
		"detail_host":    schedulerstats.ListDetail(client, schedulerstats.ListDetailOpts{HostName: "a^"}),
		"detail_backend": schedulerstats.ListDetail(client, schedulerstats.ListDetailOpts{BackendName: "a^"}),
	} {
		t.Run("no_match/"+name, func(t *testing.T) {
			pages, err := pager.AllPages(context.TODO())
			th.AssertNoErr(t, err)
			filtered, err := schedulerstats.ExtractPools(pages)
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 0, len(filtered))
		})
	}

	for name, pager := range map[string]pagination.Pager{
		"summary": schedulerstats.List(client, schedulerstats.ListOpts{ShareType: tools.RandomString("missing-", 16)}),
		"detail":  schedulerstats.ListDetail(client, schedulerstats.ListDetailOpts{ShareType: tools.RandomString("missing-", 16)}),
	} {
		t.Run("missing_share_type/"+name, func(t *testing.T) {
			_, err := pager.AllPages(context.TODO())
			th.AssertTrue(t, gophercloud.ResponseCodeIs(err, 400))
		})
	}
}
