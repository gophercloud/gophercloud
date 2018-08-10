package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/amphorae"
	fake "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestListAmphorae(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraListSuccessfully(t)

	pages := 0
	err := amphorae.List(fake.ServiceClient(), amphorae.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := amphorae.ExtractAmphoare(page)
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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraListSuccessfully(t)

	allPages, err := amphorae.List(fake.ServiceClient(), amphorae.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)
	actual, err := amphorae.ExtractAmphoare(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(actual))
}
