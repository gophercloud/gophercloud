package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListRegions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRegionsSuccessfully(t, fakeServer)

	count := 0
	err := regions.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := regions.ExtractRegions(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedRegionsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListRegionsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRegionsSuccessfully(t, fakeServer)

	allPages, err := regions.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := regions.ExtractRegions(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedRegionsSlice, actual)
	th.AssertEquals(t, ExpectedRegionsSlice[1].Extra["email"], "westsupport@example.com")
}

func TestGetRegion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetRegionSuccessfully(t, fakeServer)

	actual, err := regions.Get(context.TODO(), client.ServiceClient(fakeServer), "RegionOne-West").Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRegion, *actual)
}

func TestCreateRegion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateRegionSuccessfully(t, fakeServer)

	createOpts := regions.CreateOpts{
		ID:          "RegionOne-West",
		Description: "West sub-region of RegionOne",
		Extra: map[string]any{
			"email": "westsupport@example.com",
		},
		ParentRegionID: "RegionOne",
	}

	actual, err := regions.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRegion, *actual)
}

func TestUpdateRegion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateRegionSuccessfully(t, fakeServer)

	var description = "First West sub-region of RegionOne"
	updateOpts := regions.UpdateOpts{
		Description: &description,
		/*
			// Due to a bug in Keystone, the Extra column of the Region table
			// is not updatable, see: https://bugs.launchpad.net/keystone/+bug/1729933
			// The following lines should be uncommented once the fix is merged.

			Extra: map[string]any{
				"email": "1stwestsupport@example.com",
			},
		*/
	}

	actual, err := regions.Update(context.TODO(), client.ServiceClient(fakeServer), "RegionOne-West", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRegionUpdated, *actual)
}

func TestDeleteRegion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteRegionSuccessfully(t, fakeServer)

	res := regions.Delete(context.TODO(), client.ServiceClient(fakeServer), "RegionOne-West")
	th.AssertNoErr(t, res.Err)
}
