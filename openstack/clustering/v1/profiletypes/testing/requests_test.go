package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/clustering/v1/profiletypes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListProfileTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleList1Successfully(t)

	pageCount := 0
	err := profiletypes.List(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pageCount++
		actual, err := profiletypes.ExtractProfileTypes(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedProfileTypes, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if pageCount != 1 {
		t.Errorf("Expected 1 page, got %d", pageCount)
	}
}

func TestGetProfileType10(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGet1Successfully(t, ExpectedProfileType1.Name)

	actual, err := profiletypes.Get(context.TODO(), fake.ServiceClient(), ExpectedProfileType1.Name).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedProfileType1, *actual)
}

func TestGetProfileType15(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGet15Successfully(t, ExpectedProfileType15.Name)

	actual, err := profiletypes.Get(context.TODO(), fake.ServiceClient(), ExpectedProfileType15.Name).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedProfileType15, *actual)
}

func TestListProfileTypesOps(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListOpsSuccessfully(t)

	allPages, err := profiletypes.ListOps(fake.ServiceClient(), ProfileTypeName).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicyTypes, err := profiletypes.ExtractOps(allPages)
	th.AssertNoErr(t, err)

	for k, v := range allPolicyTypes {
		tools.PrintResource(t, k)
		tools.PrintResource(t, v)
	}
}
