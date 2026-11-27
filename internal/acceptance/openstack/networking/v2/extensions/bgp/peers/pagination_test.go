//go:build acceptance || networking || bgp || peers

package peers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBGPPeersPagination(t *testing.T) {
	clients.RequireAdmin(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	networking.RequireNeutronExtension(t, client, "bgp")
	tools.LimitListPages(t, client, client.ServiceURL("bgp-peers"))

	var expected []string
	for i := range 2 {
		peer, err := peers.Create(ctx, client, peers.CreateOpts{
			Name:     tools.RandomString("ACPTTEST-pagination-", 16),
			AuthType: "md5",
			Password: tools.MakeNewPassword(""),
			RemoteAS: 65000,
			PeerIP:   fmt.Sprintf("192.0.2.%d", i+1),
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, peers.Delete(ctx, client, peer.ID).ExtractErr())
		})
		expected = append(expected, peer.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return peers.List(client)
	}, func(page pagination.Page) ([]string, error) {
		items, err := peers.ExtractBGPPeers(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}
