package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/availabilityzoneprofiles"
	"github.com/gophercloud/gophercloud/v2/pagination"

	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListAvailabiltyZoneProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileListSuccessfully(t, fakeServer)

	pages := 0
	err := availabilityzoneprofiles.List(fake.ServiceClient(fakeServer), availabilityzoneprofiles.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := availabilityzoneprofiles.ExtractAvailabilityZoneProfiles(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 avaialbility zone profiles, got %d", len(actual))
		}
		th.CheckDeepEquals(t, AvailabilityZoneProfileSingle, actual[0])
		th.CheckDeepEquals(t, AvailabilityZoneProfileAct, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllAvailabilityZoneProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileListSuccessfully(t, fakeServer)

	allPages, err := availabilityzoneprofiles.List(fake.ServiceClient(fakeServer), availabilityzoneprofiles.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := availabilityzoneprofiles.ExtractAvailabilityZoneProfiles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AvailabilityZoneProfileSingle, actual[0])
	th.CheckDeepEquals(t, AvailabilityZoneProfileAct, actual[1])
}

func TestCreateAvailabilityZoneProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileCreationSuccessfully(t, fakeServer, SingleAvailabilityZoneProfileBody)

	actual, err := availabilityzoneprofiles.Create(context.TODO(), fake.ServiceClient(fakeServer), availabilityzoneprofiles.CreateOpts{
		Name:                 "availability-zone-profile",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"nova\"}",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, AvailabilityZoneProfileDb, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := availabilityzoneprofiles.Create(context.TODO(), fake.ServiceClient(fakeServer), availabilityzoneprofiles.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetAvailabilityZoneProfiles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := availabilityzoneprofiles.Get(context.TODO(), client, "13be083b-f502-426e-8500-07600f98b91b").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, AvailabilityZoneProfileDb, *actual)
}

func TestDeleteAvailabilityZoneProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileDeletionSuccessfully(t, fakeServer)

	res := availabilityzoneprofiles.Delete(context.TODO(), fake.ServiceClient(fakeServer), "13be083b-f502-426e-8500-07600f98b91b")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateAvailabililtyZoneProfile(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneProfileUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := availabilityzoneprofiles.Update(context.TODO(), client, "dcd65be5-f117-4260-ab3d-b32cc5bd1272", availabilityzoneprofiles.UpdateOpts{
		Name:                ptr.To("availability-zone-profile-updated"),
		ProviderName:        ptr.To("amphora"),
		AvailabiltyZoneData: ptr.To(`{"compute_zone": "nova"}`),
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, AvailabilityZoneProfileUpdated, *actual)
}
