package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/quotasets"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, "GET", "/os-quota-sets/"+FirstTenantID, getExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Get(context.TODO(), client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &getExpectedQuotaSet, actual)
}

func TestGetUsage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	uriQueryParms := map[string]string{"usage": "true"}
	HandleSuccessfulRequest(t, "GET", "/os-quota-sets/"+FirstTenantID, getUsageExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.GetUsage(context.TODO(), client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, getUsageExpectedQuotaSet, actual)
}

func TestFullUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, "PUT", "/os-quota-sets/"+FirstTenantID, fullUpdateExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, fullUpdateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &fullUpdateExpectedQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, "PUT", "/os-quota-sets/"+FirstTenantID, partialUpdateExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, partialUpdateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &partiualUpdateExpectedQuotaSet, actual)
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToBlockStorageQuotaUpdateMap() (map[string]any, error) {
	return nil, errors.New("This is an error")
}

func TestErrorInToBlockStorageQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, "PUT", "/os-quota-sets/"+FirstTenantID, "", nil)
	_, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := quotasets.Delete(context.TODO(), client.ServiceClient(), FirstTenantID).ExtractErr()
	th.AssertNoErr(t, err)
}
