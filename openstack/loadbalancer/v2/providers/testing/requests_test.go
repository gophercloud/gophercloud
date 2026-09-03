package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/providers"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListProviders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleProviderListSuccessfully(t, fakeServer)

	pages := 0
	err := providers.List(fake.ServiceClient(fakeServer), providers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := providers.ExtractProviders(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 providers, got %d", len(actual))
		}
		th.CheckDeepEquals(t, ProviderAmphora, actual[0])
		th.CheckDeepEquals(t, ProviderOVN, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllProviders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleProviderListSuccessfully(t, fakeServer)

	allPages, err := providers.List(fake.ServiceClient(fakeServer), providers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := providers.ExtractProviders(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ProviderAmphora, actual[0])
	th.CheckDeepEquals(t, ProviderOVN, actual[1])
}

func TestListFlavorCapabilities(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorCapabilitiesListSuccessfully(t, fakeServer)

	actual, err := providers.ListFlavorCapabilities(context.TODO(), fake.ServiceClient(fakeServer), "amphora", providers.ListFlavorCapabilitiesOpts{
		Fields: []string{"name", "description"},
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FlavorCapabilities, actual)
}

func TestListAvailabilityZoneCapabilities(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAvailabilityZoneCapabilitiesListSuccessfully(t, fakeServer)

	actual, err := providers.ListAvailabilityZoneCapabilities(context.TODO(), fake.ServiceClient(fakeServer), "amphora", providers.ListAvailabilityZoneCapabilitiesOpts{
		Fields: []string{"name", "description"},
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AvailabilityZoneCapabilities, actual)
}
