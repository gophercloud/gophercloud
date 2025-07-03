//go:build acceptance || dns || zone_shares

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	identity "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestShareCRD(t *testing.T) {
	// Create new project
	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := identity.CreateProject(t, identityClient, nil)
	th.AssertNoErr(t, err)
	defer identity.DeleteProject(t, identityClient, project.ID)

	// Create new Zone
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	// Create a zone share to new tenant
	share, err := CreateShare(t, client, zone, project.ID)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, share)
	defer UnshareZone(t, client, share)

	// Get the share
	getShare, err := zones.GetShare(context.TODO(), client, share.ZoneID, share.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getShare)
	th.AssertDeepEquals(t, *share, *getShare)

	// List shares
	allPages, err := zones.ListShares(client, share.ZoneID, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allShares, err := zones.ExtractZoneShares(allPages)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, allShares)

	foundShare := -1
	for i, s := range allShares {
		tools.PrintResource(t, &s)
		if share.ID == s.ID {
			foundShare = i
			break
		}
	}
	if foundShare == -1 {
		t.Fatalf("Share %s not found in list", share.ID)
	}

	th.AssertDeepEquals(t, *share, allShares[foundShare])
}
