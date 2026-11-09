//go:build acceptance || compute || servers || instanceactions || migrations || secgroups || volumeattach

package v2

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// createPaginationServer registers cleanup before waiting for the server to
// become active. It returns the create request ID as an independent action ID.
func createPaginationServer(t *testing.T, ctx context.Context, client *gophercloud.ServiceClient, groups ...string) (*servers.Server, string) {
	t.Helper()
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)
	networkID := choices.NetworkID
	if networkID == "" {
		networkClient, err := clients.NewNetworkV2Client()
		th.AssertNoErr(t, err)
		page, err := networks.List(networkClient, networks.ListOpts{Name: choices.NetworkName}).AllPages(ctx)
		th.AssertNoErr(t, err)
		items, err := networks.ExtractNetworks(page)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(items))
		networkID = items[0].ID
	}
	result := servers.Create(ctx, client, servers.CreateOpts{
		Name:           tools.RandomString("ACPTTEST-pagination-", 16),
		FlavorRef:      choices.FlavorID,
		ImageRef:       choices.ImageID,
		Networks:       []servers.Network{{UUID: networkID}},
		SecurityGroups: groups,
	}, nil)
	server, err := result.Extract()
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		err := servers.Delete(ctx, client, server.ID).ExtractErr()
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return
		}
		if err != nil {
			t.Errorf("Unable to delete server %s: %v", server.ID, err)
			return
		}
		th.CheckNoErr(t, gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
			_, err := servers.Get(ctx, client, server.ID).Extract()
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return true, nil
			}
			return false, err
		}))
	})
	th.AssertNoErr(t, waitPaginationServerStatus(ctx, client, server.ID, "ACTIVE"))
	return server, result.Header.Get("X-OpenStack-Request-ID")
}

func waitPaginationServerStatus(ctx context.Context, client *gophercloud.ServiceClient, serverID, status string) error {
	return gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
		server, err := servers.Get(ctx, client, serverID).Extract()
		if err != nil {
			return false, err
		}
		if server.Status == "ERROR" {
			return false, fmt.Errorf("server %s entered ERROR while waiting for %s", serverID, status)
		}
		return server.Status == status, nil
	})
}
