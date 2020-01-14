package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListResourceProviders(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleResourceProviderList(t)

	count := 0
	err := resourceproviders.List(fake.ServiceClient(), resourceproviders.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := resourceproviders.ExtractResourceProviders(page)
		if err != nil {
			t.Errorf("Failed to extract resource providers: %v", err)
			return false, err
		}
		th.AssertDeepEquals(t, ExpectedResourceProviders, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
