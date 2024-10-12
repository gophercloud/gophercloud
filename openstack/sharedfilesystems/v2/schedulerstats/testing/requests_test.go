package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/schedulerstats"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListPoolsDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePoolsListSuccessfully(t)

	pages := 0
	err := schedulerstats.List(client.ServiceClient(), schedulerstats.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := schedulerstats.ExtractPools(page)
		th.AssertNoErr(t, err)

		if len(actual) != 4 {
			t.Fatalf("Expected 4 backends, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PoolFake1, actual[0])
		th.CheckDeepEquals(t, PoolFake2, actual[1])
		th.CheckDeepEquals(t, PoolFake3, actual[2])
		th.CheckDeepEquals(t, PoolFake4, actual[3])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}

	pages = 0
	err = schedulerstats.ListDetail(client.ServiceClient(), schedulerstats.ListDetailOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := schedulerstats.ExtractPools(page)
		th.AssertNoErr(t, err)

		if len(actual) != 4 {
			t.Fatalf("Expected 4 backends, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PoolDetailFake1, actual[0])
		th.CheckDeepEquals(t, PoolDetailFake2, actual[1])
		th.CheckDeepEquals(t, PoolDetailFake3, actual[2])
		th.CheckDeepEquals(t, PoolDetailFake4, actual[3])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
