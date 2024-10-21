//go:build acceptance || identity || regions

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/regions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRegionsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := regions.ListOpts{
		ParentRegionID: "RegionOne",
	}

	allPages, err := regions.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRegions, err := regions.ExtractRegions(allPages)
	th.AssertNoErr(t, err)

	for _, region := range allRegions {
		tools.PrintResource(t, region)
	}
}

func TestRegionsGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := regions.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRegions, err := regions.ExtractRegions(allPages)
	th.AssertNoErr(t, err)

	region := allRegions[0]
	p, err := regions.Get(context.TODO(), client, region.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, p)

	th.AssertEquals(t, region.ID, p.ID)
}

func TestRegionsCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := regions.CreateOpts{
		ID:          "testregion",
		Description: "Region for testing",
		Extra: map[string]any{
			"email": "testregion@example.com",
		},
	}

	// Create region in the default domain
	region, err := CreateRegion(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteRegion(t, client, region.ID)

	tools.PrintResource(t, region)
	tools.PrintResource(t, region.Extra)

	var description = ""
	updateOpts := regions.UpdateOpts{
		Description: &description,
		/*
			// Due to a bug in Keystone, the Extra column of the Region table
			// is not updatable, see: https://bugs.launchpad.net/keystone/+bug/1729933
			// The following lines should be uncommented once the fix is merged.

			Extra: map[string]any{
				"email": "testregionA@example.com",
			},
		*/
	}

	newRegion, err := regions.Update(context.TODO(), client, region.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newRegion)
	tools.PrintResource(t, newRegion.Extra)

	th.AssertEquals(t, newRegion.Description, description)
}
