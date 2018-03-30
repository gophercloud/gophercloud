package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t, "/os-quota-sets/"+FirstTenantID, GetOutput)
	actual, err := quotasets.Get(client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestGetDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t, "/os-quota-sets/"+FirstTenantID+"/detail", GetDetailsOutput)
	actual, err := quotasets.GetDetail(client.ServiceClient(), FirstTenantID).Extract()
	th.CheckDeepEquals(t, FirstQuotaDetailsSet, actual)
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t, "/os-quota-sets/"+FirstTenantID, UpdateOutput)
	actual, err := quotasets.Update(client.ServiceClient(), FirstTenantID, UpdatedQuotaSet).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstQuotaSet, actual)
}

func TestPartialUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t, "/os-quota-sets/"+FirstTenantID, PartialUpdateBody)
	opts := quotasets.UpdateOpts{Volumes: gophercloud.IntToPointer(200), Force: true}
	actual, err := quotasets.Update(client.ServiceClient(), FirstTenantID, opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &PartialQuotaSet, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)
	_, err := quotasets.Delete(client.ServiceClient(), FirstTenantID).Extract()
	th.AssertNoErr(t, err)
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToBlockStorageQuotaUpdateMap() (map[string]interface{}, error) {
	return nil, errors.New("This is an error")
}

func TestErrorInToBlockStorageQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePutSuccessfully(t, "/os-quota-sets/"+FirstTenantID, UpdateOutput)
	_, err := quotasets.Update(client.ServiceClient(), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}
