//go:build acceptance || instanceha

package v1

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/hosts"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/notifications"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/segments"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/vmoves"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestInstanceHAReadOnly(t *testing.T) {
	clients.RequireAdmin(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	opts, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)
	provider, err := openstack.AuthenticatedClient(ctx, opts)
	th.AssertNoErr(t, err)
	client, err := openstack.NewInstanceHAV1(ctx, provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	th.AssertNoErr(t, err)
	client.Microversion = "1.3"
	segmentPages, err := segments.List(client, segments.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
	th.AssertNoErr(t, err)
	segmentItems, err := segments.ExtractSegments(segmentPages)
	th.AssertNoErr(t, err)
	for _, segment := range segmentItems {
		_, err := segments.Get(ctx, client, segment.UUID).Extract()
		th.AssertNoErr(t, err)
		pages, err := hosts.List(client, segment.UUID, hosts.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
		th.AssertNoErr(t, err)
		items, err := hosts.ExtractHosts(pages)
		th.AssertNoErr(t, err)
		for _, host := range items {
			_, err := hosts.Get(ctx, client, segment.UUID, host.UUID).Extract()
			th.AssertNoErr(t, err)
		}
	}
	pages, err := notifications.List(client, notifications.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc", GeneratedSince: time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)}).AllPages(ctx)
	th.AssertNoErr(t, err)
	items, err := notifications.ExtractNotifications(pages)
	th.AssertNoErr(t, err)
	for _, notification := range items {
		_, err := notifications.Get(ctx, client, notification.UUID).Extract()
		th.AssertNoErr(t, err)
		if notification.Type != "COMPUTE_HOST" {
			continue
		}
		pages, err := vmoves.List(client, notification.UUID, vmoves.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
		th.AssertNoErr(t, err)
		moves, err := vmoves.ExtractVMoves(pages)
		th.AssertNoErr(t, err)
		for _, move := range moves {
			_, err := vmoves.Get(ctx, client, notification.UUID, move.UUID).Extract()
			th.AssertNoErr(t, err)
		}
	}
	t.Logf("Read %d segments and %d recent notifications; empty clouds cannot exercise resource detail endpoints", len(segmentItems), len(items))
}
