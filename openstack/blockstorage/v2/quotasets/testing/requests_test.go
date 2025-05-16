package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/quotasets"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, fakeServer, "GET", "/os-quota-sets/"+FirstTenantID, getExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Get(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &getExpectedQuotaSet, actual)
}

func TestGetUsage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	uriQueryParms := map[string]string{"usage": "true"}
	HandleSuccessfulRequest(t, fakeServer, "GET", "/os-quota-sets/"+FirstTenantID, getUsageExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.GetUsage(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, getUsageExpectedQuotaSet, actual)
}

func TestFullUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, fakeServer, "PUT", "/os-quota-sets/"+FirstTenantID, fullUpdateExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, fullUpdateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &fullUpdateExpectedQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	uriQueryParms := map[string]string{}
	HandleSuccessfulRequest(t, fakeServer, "PUT", "/os-quota-sets/"+FirstTenantID, partialUpdateExpectedJSONBody, uriQueryParms)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, partialUpdateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &partiualUpdateExpectedQuotaSet, actual)
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToBlockStorageQuotaUpdateMap() (map[string]any, error) {
	return nil, errors.New("this is an error")
}

func TestErrorInToBlockStorageQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSuccessfulRequest(t, fakeServer, "PUT", "/os-quota-sets/"+FirstTenantID, "", nil)
	_, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	err := quotasets.Delete(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).ExtractErr()
	th.AssertNoErr(t, err)
}
