package accounts

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/objectstorage/v1/accounts"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetAccounts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleGetAccountSuccessfully(t)

	options := &UpdateOpts{Metadata: map[string]string{"gophercloud-test": "accounts"}}
	_, err := Update(fake.ServiceClient(), options).ExtractHeaders()
	if err != nil {
		t.Fatalf("Unable to update account: %v", err)
	}
}

func TestUpdateAccounts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleUpdateAccountSuccessfully(t)

	expected := map[string]string{"Foo": "bar"}
	actual, err := Get(fake.ServiceClient()).ExtractMetadata()
	if err != nil {
		t.Fatalf("Unable to get account metadata: %v", err)
	}
	th.CheckDeepEquals(t, expected, actual)
}
