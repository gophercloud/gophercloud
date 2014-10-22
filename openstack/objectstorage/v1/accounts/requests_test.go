package accounts

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "accounts"}

func TestUpdateAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAccountSuccessfully(t)

	options := &UpdateOpts{Metadata: map[string]string{"gophercloud-test": "accounts"}}
	_, err := Update(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestGetAccount(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateAccountSuccessfully(t)

	expected := map[string]string{"Foo": "bar"}
	actual, err := Get(fake.ServiceClient(), &GetOpts{}).ExtractMetadata()
	if err != nil {
		t.Fatalf("Unable to get account metadata: %v", err)
	}
	th.CheckDeepEquals(t, expected, actual)
}
