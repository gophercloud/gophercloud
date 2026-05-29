package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)
	actual, err := quotasets.Get(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestGetDetail(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetDetailSuccessfully(t, fakeServer)
	actual, err := quotasets.GetDetail(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).Extract()
	th.CheckDeepEquals(t, FirstQuotaDetailsSet, actual)
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePutSuccessfully(t, fakeServer)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, UpdatedQuotaSet).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePartialPutSuccessfully(t, fakeServer)
	opts := quotasets.UpdateOpts{Cores: gophercloud.IntToPointer(200), Force: true}
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)
	_, err := quotasets.Delete(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToComputeQuotaUpdateMap() (map[string]any, error) {
	return nil, errors.New("this is an error")
}

func TestErrorInToComputeQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePutSuccessfully(t, fakeServer)
	_, err := quotasets.Update(context.TODO(), client.ServiceClient(fakeServer), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}
