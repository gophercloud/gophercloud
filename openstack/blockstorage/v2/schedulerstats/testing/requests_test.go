package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/schedulerstats"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListStoragePoolsDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleStoragePoolsListSuccessfully(t)

	pages := 0
	err := schedulerstats.List(client.ServiceClient(), schedulerstats.ListOpts{Detail: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := schedulerstats.ExtractStoragePools(page)
		th.AssertNoErr(t, err)

		if len(actual) != 2 {
			t.Fatalf("Expected 2 backends, got %d", len(actual))
		}
		th.CheckDeepEquals(t, StoragePoolFake1, actual[0])
		th.CheckDeepEquals(t, StoragePoolFake2, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
