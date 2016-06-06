package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/accounts"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestUpdateAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateAccountSuccessfully(t)

	options := &accounts.UpdateOpts{Metadata: map[string]string{"gophercloud-test": "accounts"}}
	_, err := accounts.Update(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestGetAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAccountSuccessfully(t)

	expectedMetadata := map[string]string{"Subject": "books"}
	res := accounts.Get(fake.ServiceClient(), &accounts.GetOpts{})
	th.AssertNoErr(t, res.Err)
	actualMetadata, _ := res.ExtractMetadata()
	th.CheckDeepEquals(t, expectedMetadata, actualMetadata)
	_, err := res.Extract()
	th.AssertNoErr(t, err)
}
