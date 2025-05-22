package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/usage"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetTenant(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSingleTenantSuccessfully(t, fakeServer)

	count := 0
	err := usage.SingleTenant(client.ServiceClient(fakeServer), FirstTenantID, nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := usage.ExtractSingleTenant(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, &SingleTenantUsageResults, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestAllTenants(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetAllTenantsSuccessfully(t, fakeServer)

	getOpts := usage.AllTenantsOpts{
		Detailed: true,
	}

	count := 0
	err := usage.AllTenants(client.ServiceClient(fakeServer), getOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := usage.ExtractAllTenants(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, AllTenantsUsageResult, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}
