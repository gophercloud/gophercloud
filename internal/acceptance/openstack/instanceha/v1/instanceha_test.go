//go:build acceptance || instanceha

package v1

import (
	"context"
	"os"
	"slices"
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
	hostCount, moveCount := 0, 0
	segmentPages, err := segments.List(client, segments.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
	th.AssertNoErr(t, err)
	segmentItems, err := segments.ExtractSegments(segmentPages)
	th.AssertNoErr(t, err)
	for segmentIndex, segment := range segmentItems {
		if segmentIndex == 0 {
			pages, err := segments.List(client, segments.ListOpts{Limit: 1, SortKey: "id", SortDir: "desc", RecoveryMethod: segment.RecoveryMethod, ServiceType: segment.ServiceType, Enabled: segment.Enabled}).AllPages(ctx)
			th.AssertNoErr(t, err)
			filtered, err := segments.ExtractSegments(pages)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, true, slices.ContainsFunc(filtered, func(item segments.Segment) bool { return item.UUID == segment.UUID }))
			for _, item := range filtered {
				th.AssertEquals(t, segment.RecoveryMethod, item.RecoveryMethod)
				th.AssertEquals(t, segment.ServiceType, item.ServiceType)
				if segment.Enabled != nil {
					if item.Enabled == nil {
						t.Fatal("filtered segment is missing enabled")
					}
					th.AssertEquals(t, *segment.Enabled, *item.Enabled)
				}
			}
		}
		detail, err := segments.Get(ctx, client, segment.UUID).Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, segment.UUID, detail.UUID)
		pages, err := hosts.List(client, segment.UUID, hosts.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
		th.AssertNoErr(t, err)
		items, err := hosts.ExtractHosts(pages)
		th.AssertNoErr(t, err)
		hostCount += len(items)
		for hostIndex, host := range items {
			if hostIndex == 0 {
				pages, err := hosts.List(client, segment.UUID, hosts.ListOpts{Limit: 1, SortKey: "id", SortDir: "desc", Name: host.Name, Type: host.Type, ControlAttributes: host.ControlAttributes, Reserved: &host.Reserved, OnMaintenance: &host.OnMaintenance}).AllPages(ctx)
				th.AssertNoErr(t, err)
				filtered, err := hosts.ExtractHosts(pages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, true, slices.ContainsFunc(filtered, func(item hosts.Host) bool { return item.UUID == host.UUID }))
				for _, item := range filtered {
					th.AssertEquals(t, host.Name, item.Name)
					th.AssertEquals(t, host.Type, item.Type)
					th.AssertEquals(t, host.ControlAttributes, item.ControlAttributes)
					th.AssertEquals(t, host.Reserved, item.Reserved)
					th.AssertEquals(t, host.OnMaintenance, item.OnMaintenance)
				}
			}
			detail, err := hosts.Get(ctx, client, segment.UUID, host.UUID).Extract()
			th.AssertNoErr(t, err)
			th.AssertEquals(t, host.UUID, detail.UUID)
			th.AssertEquals(t, segment.UUID, detail.FailoverSegmentID)
			if detail.FailoverSegment == nil {
				t.Fatal("host response is missing its embedded failover segment")
			}
			th.AssertEquals(t, segment.UUID, detail.FailoverSegment.UUID)
			if time.Time(detail.CreatedAt).IsZero() || time.Time(detail.FailoverSegment.CreatedAt).IsZero() {
				t.Fatal("host or embedded segment creation time is missing")
			}
		}
	}
	pages, err := notifications.List(client, notifications.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc", GeneratedSince: time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)}).AllPages(ctx)
	th.AssertNoErr(t, err)
	items, err := notifications.ExtractNotifications(pages)
	th.AssertNoErr(t, err)
	for notificationIndex, notification := range items {
		if notificationIndex == 0 {
			pages, err := notifications.List(client, notifications.ListOpts{Limit: 1, SortKey: "id", SortDir: "desc", Status: notification.Status, Type: notification.Type, SourceHostUUID: notification.SourceHostUUID, GeneratedSince: time.Time(notification.GeneratedTime).Add(-time.Second).Format(time.RFC3339)}).AllPages(ctx)
			th.AssertNoErr(t, err)
			filtered, err := notifications.ExtractNotifications(pages)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, true, slices.ContainsFunc(filtered, func(item notifications.Notification) bool { return item.UUID == notification.UUID }))
			for _, item := range filtered {
				th.AssertEquals(t, notification.Status, item.Status)
				th.AssertEquals(t, notification.Type, item.Type)
				th.AssertEquals(t, notification.SourceHostUUID, item.SourceHostUUID)
			}
		}
		detail, err := notifications.Get(ctx, client, notification.UUID).Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, notification.UUID, detail.UUID)
		if notification.Type != "COMPUTE_HOST" {
			continue
		}
		pages, err := vmoves.List(client, notification.UUID, vmoves.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc"}).AllPages(ctx)
		th.AssertNoErr(t, err)
		moves, err := vmoves.ExtractVMoves(pages)
		th.AssertNoErr(t, err)
		moveCount += len(moves)
		for moveIndex, move := range moves {
			if moveIndex == 0 {
				pages, err := vmoves.List(client, notification.UUID, vmoves.ListOpts{Limit: 1, SortKey: "id", SortDir: "desc", Status: move.Status, Type: move.Type}).AllPages(ctx)
				th.AssertNoErr(t, err)
				filtered, err := vmoves.ExtractVMoves(pages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, true, slices.ContainsFunc(filtered, func(item vmoves.VMove) bool { return item.UUID == move.UUID }))
				for _, item := range filtered {
					th.AssertEquals(t, move.Status, item.Status)
					th.AssertEquals(t, move.Type, item.Type)
				}
			}
			detail, err := vmoves.Get(ctx, client, notification.UUID, move.UUID).Extract()
			th.AssertNoErr(t, err)
			th.AssertEquals(t, move.UUID, detail.UUID)
		}
	}
	t.Logf("Read %d segments, %d hosts, %d recent notifications and %d VM moves", len(segmentItems), hostCount, len(items), moveCount)
}
