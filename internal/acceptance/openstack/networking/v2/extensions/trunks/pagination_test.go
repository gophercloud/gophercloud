//go:build acceptance || networking || trunks

package trunks

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// trunkPaginationOpts adds the page size missing from trunks.ListOpts.
type trunkPaginationOpts struct{}

func (trunkPaginationOpts) ToTrunkListQuery() (string, error) {
	return "?limit=1&sort_key=id&sort_dir=asc", nil
}

func TestTrunksPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	networking.RequireNeutronExtension(t, client, "trunk")

	network, err := networks.Create(ctx, client, networks.CreateOpts{
		Name: tools.RandomString("ACPTTEST-trunk-pagination-", 16),
	}).Extract()
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		th.CheckNoErr(t, networks.Delete(ctx, client, network.ID).ExtractErr())
	})

	var expected []string
	for range 2 {
		port, err := ports.Create(ctx, client, ports.CreateOpts{NetworkID: network.ID}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, ports.Delete(ctx, client, port.ID).ExtractErr())
		})
		trunk, err := trunks.Create(ctx, client, trunks.CreateOpts{
			Name:   tools.RandomString("ACPTTEST-pagination-", 16),
			PortID: port.ID,
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, trunks.Delete(ctx, client, trunk.ID).ExtractErr())
		})
		expected = append(expected, trunk.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return trunks.List(client, trunkPaginationOpts{})
	}, func(page pagination.Page) ([]string, error) {
		items, err := trunks.ExtractTrunks(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}
