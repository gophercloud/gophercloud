//go:build acceptance || networking || networkipavailabilities

package networkipavailabilities

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/networkipavailabilities"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// availabilityPaginationOpts adds page size and sorting to the list query.
type availabilityPaginationOpts struct{}

func (availabilityPaginationOpts) ToNetworkIPAvailabilityListQuery() (string, error) {
	return "?limit=1&sort_key=network_id&sort_dir=asc", nil
}

func TestNetworkIPAvailabilityPagination(t *testing.T) {
	clients.RequireAdmin(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	networking.RequireNeutronExtension(t, client, "network-ip-availability")

	var expected []string
	for range 2 {
		network, err := networks.Create(ctx, client, networks.CreateOpts{
			Name: tools.RandomString("ACPTTEST-ip-availability-pagination-", 16),
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, networks.Delete(ctx, client, network.ID).ExtractErr())
		})
		dhcp := false
		gateway := ""
		subnet, err := subnets.Create(ctx, client, subnets.CreateOpts{
			NetworkID:  network.ID,
			CIDR:       "192.0.2.0/24",
			IPVersion:  gophercloud.IPv4,
			EnableDHCP: &dhcp,
			GatewayIP:  &gateway,
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, subnets.Delete(ctx, client, subnet.ID).ExtractErr())
		})
		expected = append(expected, network.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return networkipavailabilities.List(client, availabilityPaginationOpts{})
	}, func(page pagination.Page) ([]string, error) {
		items, err := networkipavailabilities.ExtractNetworkIPAvailabilities(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.NetworkID
		}
		return ids, err
	}, expected)
}
