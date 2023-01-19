package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack/identity/v3/extensions/federation"
	"github.com/bizflycloud/gophercloud/pagination"
	th "github.com/bizflycloud/gophercloud/testhelper"
	"github.com/bizflycloud/gophercloud/testhelper/client"
)

func TestListMappings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListMappingsSuccessfully(t)

	count := 0
	err := federation.ListMappings(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := federation.ExtractMappings(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedMappingsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListMappingsAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListMappingsSuccessfully(t)

	allPages, err := federation.ListMappings(client.ServiceClient()).AllPages()
	th.AssertNoErr(t, err)
	actual, err := federation.ExtractMappings(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedMappingsSlice, actual)
}
