package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetTenant(t *testing.T) {
	var getOpts usage.SingleTenantOpts
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSingleTenantSuccessfully(t)
	page, err := usage.SingleTenant(client.ServiceClient(), FirstTenantID, getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := usage.ExtractSingleTenant(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SingleTenantUsageResults, actual)
}

func TestAllTenants(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAllTenantsSuccessfully(t)

	getOpts := usage.AllTenantsOpts{
		Detailed: true,
	}

	page, err := usage.AllTenants(client.ServiceClient(), getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := usage.ExtractAllTenants(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AllTenantsUsageResult, actual)
}
