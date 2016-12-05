package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/gophercloud/gophercloud"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)
	actual, err := quotasets.Get(client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t)
	actual, err := quotasets.Update(client.ServiceClient(), FirstTenantID, UpdatedQuotaSet).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePartialPutSuccessfully(t)
	opts := quotasets.UpdateOpts{Cores:gophercloud.IntToPointer(200), Force:true}
	actual, err := quotasets.Update(client.ServiceClient(), FirstTenantID, opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)
	_, err := quotasets.Delete(client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
}