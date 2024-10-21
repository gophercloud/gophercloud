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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)
	actual, err := quotasets.Get(context.TODO(), client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestGetDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetDetailSuccessfully(t)
	actual, err := quotasets.GetDetail(context.TODO(), client.ServiceClient(), FirstTenantID).Extract()
	th.CheckDeepEquals(t, FirstQuotaDetailsSet, actual)
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t)
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, UpdatedQuotaSet).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePartialPutSuccessfully(t)
	opts := quotasets.UpdateOpts{Cores: gophercloud.IntToPointer(200), Force: true}
	actual, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)
	_, err := quotasets.Delete(context.TODO(), client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToComputeQuotaUpdateMap() (map[string]any, error) {
	return nil, errors.New("This is an error")
}

func TestErrorInToComputeQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t)
	_, err := quotasets.Update(context.TODO(), client.ServiceClient(), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}
