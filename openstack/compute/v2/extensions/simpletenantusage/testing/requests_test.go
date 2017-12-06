package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/simpletenantusage"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetTenant(t *testing.T) {
	var getOpts simpletenantusage.SingleTenantOpts
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSingleTenantSuccessfully(t)
	page, err := simpletenantusage.SingleTenant(client.ServiceClient(), FirstTenantID, getOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := simpletenantusage.ExtractSingleTenant(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SingleTenantUsageResults, actual)
}
