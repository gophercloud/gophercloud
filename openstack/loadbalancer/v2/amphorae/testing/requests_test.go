package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/amphorae"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListAmphorae(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAmphoraListSuccessfully(t, fakeServer)

	pages := 0
	err := amphorae.List(fake.ServiceClient(fakeServer), amphorae.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := amphorae.ExtractAmphorae(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 amphorae, got %d", len(actual))
		}

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllAmphorae(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAmphoraListSuccessfully(t, fakeServer)

	allPages, err := amphorae.List(fake.ServiceClient(fakeServer), amphorae.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := amphorae.ExtractAmphorae(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(actual))
	th.AssertDeepEquals(t, ExpectedAmphoraeSlice, actual)
}

func TestGetAmphora(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAmphoraGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := amphorae.Get(context.TODO(), client, "45f40289-0551-483a-b089-47214bc2a8a4").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, FirstAmphora, *actual)
}

func TestFailoverAmphora(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAmphoraFailoverSuccessfully(t, fakeServer)

	res := amphorae.Failover(context.TODO(), fake.ServiceClient(fakeServer), "36e08a3e-a78f-4b40-a229-1e7e23eee1ab")
	th.AssertNoErr(t, res.Err)
}
