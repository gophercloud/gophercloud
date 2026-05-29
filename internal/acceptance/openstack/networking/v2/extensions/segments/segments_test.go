//go:build acceptance || networking || segments

package segments

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSegmentCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "segment")

	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	segment, err := CreateSegment(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSegment(t, client, segment.ID)

	// Get
	segGet, err := segments.Get(context.TODO(), client, segment.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, segment.ID, segGet.ID)

	// Update
	newName := tools.RandomString("UPDATED-SEGMENT-", 8)
	newDesc := "updated description"
	updateOpts := segments.UpdateOpts{
		Name:        &newName,
		Description: &newDesc,
	}
	segUpdated, err := segments.Update(context.TODO(), client, segment.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newName, segUpdated.Name)
	th.AssertEquals(t, newDesc, segUpdated.Description)

	// List
	allPages, err := segments.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allSegments, err := segments.ExtractSegments(allPages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(allSegments), 1)
	t.Logf("Found %d segments", len(allSegments))
	tools.PrintResource(t, allSegments)
}
