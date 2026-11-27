//go:build acceptance || compute || migrations || servers

package v2

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/migrations"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestMigrationsPagination(t *testing.T) {
	clients.RequireLong(t)
	clients.RequireAdmin(t)
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.59"
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	server, _ := createPaginationServer(t, ctx, client)

	// Revert each resize so the same base/resize flavors can be reused. This
	// does not require the optional OS_LIVE_MIGRATE environment setting.
	for cycle := 0; cycle < 2; cycle++ {
		th.AssertNoErr(t, servers.Resize(ctx, client, server.ID, servers.ResizeOpts{FlavorRef: choices.FlavorIDResize}).ExtractErr())
		th.AssertNoErr(t, waitPaginationServerStatus(ctx, client, server.ID, "VERIFY_RESIZE"))
		th.AssertNoErr(t, servers.RevertResize(ctx, client, server.ID).ExtractErr())
		th.AssertNoErr(t, waitPaginationServerStatus(ctx, client, server.ID, "ACTIVE"))
		th.AssertNoErr(t, gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
			items := migrationPaginationControl(t, ctx, client, server.ID)
			th.AssertEquals(t, cycle+1, len(items))
			for _, item := range items {
				if item.Status == "error" {
					return false, fmt.Errorf("migration %s entered error", item.UUID)
				}
				if item.Status != "reverted" {
					return false, nil
				}
			}
			return true, nil
		}))
	}

	// Independent, explicit-marker control GETs establish the two identities.
	// They do not use the SDK pager or the API's continuation links.
	items := migrationPaginationControl(t, ctx, client, server.ID)
	th.AssertEquals(t, 2, len(items))
	expected := []string{items[0].UUID, items[1].UUID}
	limit := 1
	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return migrations.List(client, migrations.ListOpts{InstanceID: &server.ID, Limit: &limit})
	}, func(page pagination.Page) ([]string, error) {
		items, err := migrations.ExtractMigrations(page)
		ids := make([]string, len(items))
		for i, item := range items {
			if item.InstanceID != server.ID {
				return nil, fmt.Errorf("migration belongs to server %s, expected %s", item.InstanceID, server.ID)
			}
			ids[i] = item.UUID
		}
		return ids, err
	}, expected)
}

func migrationPaginationControl(t *testing.T, ctx context.Context, client *gophercloud.ServiceClient, serverID string) []migrations.Migration {
	t.Helper()
	var items []migrations.Migration
	seen := make(map[string]bool)
	query := url.Values{"instance_uuid": {serverID}, "limit": {"1"}}
	for {
		var body struct {
			Migrations []migrations.Migration `json:"migrations"`
		}
		_, err := client.Get(ctx, client.ServiceURL("os-migrations")+"?"+query.Encode(), &body, nil)
		th.AssertNoErr(t, err)
		if len(body.Migrations) == 0 {
			return items
		}
		th.AssertEquals(t, 1, len(body.Migrations))
		item := body.Migrations[0]
		th.AssertEquals(t, serverID, item.InstanceID)
		th.AssertEquals(t, "resize", item.MigrationType)
		if item.UUID == "" || seen[item.UUID] || len(items) >= 2 {
			t.Fatalf("unexpected migration control result: %+v", item)
		}
		seen[item.UUID] = true
		items = append(items, item)
		query.Set("marker", item.UUID)
	}
}
