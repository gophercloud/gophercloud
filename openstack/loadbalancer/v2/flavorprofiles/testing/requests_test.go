package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
	"github.com/gophercloud/gophercloud/v2/pagination"

	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListFlavorProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileListSuccessfully(t, fakeServer)

	pages := 0
	err := flavorprofiles.List(fake.ServiceClient(fakeServer), flavorprofiles.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := flavorprofiles.ExtractFlavorProfiles(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 flavors, got %d", len(actual))
		}
		th.CheckDeepEquals(t, FlavorProfileSingle, actual[0])
		th.CheckDeepEquals(t, FlavorProfileAct, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllFlavorProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileListSuccessfully(t, fakeServer)

	allPages, err := flavorprofiles.List(fake.ServiceClient(fakeServer), flavorprofiles.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := flavorprofiles.ExtractFlavorProfiles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FlavorProfileSingle, actual[0])
	th.CheckDeepEquals(t, FlavorProfileAct, actual[1])
}

func TestCreateFlavorProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileCreationSuccessfully(t, fakeServer, SingleFlavorProfileBody)

	actual, err := flavorprofiles.Create(context.TODO(), fake.ServiceClient(fakeServer), flavorprofiles.CreateOpts{
		Name:         "amphora-test",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := flavorprofiles.Create(context.TODO(), fake.ServiceClient(fakeServer), flavorprofiles.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetFlavorProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := flavorprofiles.Get(context.TODO(), client, "dcd65be5-f117-4260-ab3d-b32cc5bd1272").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestDeleteFlavorProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileDeletionSuccessfully(t, fakeServer)

	res := flavorprofiles.Delete(context.TODO(), fake.ServiceClient(fakeServer), "dcd65be5-f117-4260-ab3d-b32cc5bd1272")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateFlavorProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorProfileUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := flavorprofiles.Update(context.TODO(), client, "dcd65be5-f117-4260-ab3d-b32cc5bd1272", flavorprofiles.UpdateOpts{
		Name:         "amphora-test-updated",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"SINGLE\"}",
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorUpdated, *actual)
}
