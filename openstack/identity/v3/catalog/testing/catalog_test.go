package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack/identity/v3/catalog"
	"github.com/bizflycloud/gophercloud/pagination"
	th "github.com/bizflycloud/gophercloud/testhelper"
	"github.com/bizflycloud/gophercloud/testhelper/client"
)

func TestListCatalog(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListCatalogSuccessfully(t)

	count := 0
	err := catalog.List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := catalog.ExtractServiceCatalog(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedCatalogSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}
