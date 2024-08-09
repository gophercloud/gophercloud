package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/services"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListServices(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	pages := 0
	err := services.List(client.ServiceClient(), services.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 3 {
			t.Fatalf("Expected 3 services, got %d", len(actual))
		}
		th.CheckDeepEquals(t, FirstFakeService, actual[0])
		th.CheckDeepEquals(t, SecondFakeService, actual[1])
		th.CheckDeepEquals(t, ThirdFakeService, actual[2])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
