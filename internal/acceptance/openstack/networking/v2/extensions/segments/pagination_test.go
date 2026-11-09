//go:build acceptance || networking || segments

package segments

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// segmentPaginationOpts adds the page size missing from segments.ListOpts.
type segmentPaginationOpts struct{}

func (segmentPaginationOpts) ToSegmentListQuery() (string, error) {
	return "?limit=1&sort_key=id&sort_dir=asc", nil
}

func TestSegmentsPagination(t *testing.T) {
	clients.RequireAdmin(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	networking.RequireNeutronExtension(t, client, "segment")

	var expected []string
	for range 2 {
		network, err := networks.Create(ctx, client, networks.CreateOpts{
			Name: tools.RandomString("ACPTTEST-segment-pagination-", 16),
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, networks.Delete(ctx, client, network.ID).ExtractErr())
		})
		segment, err := segments.Create(ctx, client, segments.CreateOpts{
			Name:        tools.RandomString("ACPTTEST-pagination-", 16),
			NetworkID:   network.ID,
			NetworkType: "geneve",
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, segments.Delete(ctx, client, segment.ID).ExtractErr())
		})
		expected = append(expected, segment.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return segments.List(client, segmentPaginationOpts{})
	}, func(page pagination.Page) ([]string, error) {
		items, err := segments.ExtractSegments(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}
