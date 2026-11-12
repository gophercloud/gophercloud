//go:build acceptance || networking || bgp || speakers

package speakers

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBGPSpeakersPagination(t *testing.T) {
	clients.RequireAdmin(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	networking.RequireNeutronExtension(t, client, "bgp")
	tools.LimitListPages(t, client, client.ServiceURL("bgp-speakers"))

	var expected []string
	for range 2 {
		speaker, err := speakers.Create(ctx, client, speakers.CreateOpts{
			Name:      tools.RandomString("ACPTTEST-pagination-", 16),
			IPVersion: 4,
			LocalAS:   65000,
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, speakers.Delete(ctx, client, speaker.ID).ExtractErr())
		})
		expected = append(expected, speaker.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return speakers.List(client)
	}, func(page pagination.Page) ([]string, error) {
		items, err := speakers.ExtractBGPSpeakers(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}
