package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/simpletenantusage"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGet(t *testing.T) {
	var getOpts simpletenantusage.GetOpts
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)
	page, err := simpletenantusage.Get(client.ServiceClient(), getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := simpletenantusage.ExtractSimpleTenantUsages(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SimpleTenantUsageResults, actual)
}

func TestGetTenant(t *testing.T) {
	var getOpts simpletenantusage.GetOpts
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTenantSuccessfully(t)
	page, err := simpletenantusage.GetTenant(client.ServiceClient(), FirstTenantID, getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := simpletenantusage.ExtractSimpleTenantUsage(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SimpleTenantUsageOneTenantResults, actual)
}

func TestGetDetails(t *testing.T) {
	detailed := true
	getOpts := simpletenantusage.GetOpts{
		Detailed: &detailed,
	}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDetailSuccessfully(t)
	page, err := simpletenantusage.Get(client.ServiceClient(), getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := simpletenantusage.ExtractSimpleTenantUsages(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SimpleTenantUsageDetailResults, actual)
}
