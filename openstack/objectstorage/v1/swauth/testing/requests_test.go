package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack"
	"github.com/bizflycloud/gophercloud/openstack/objectstorage/v1/swauth"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestAuth(t *testing.T) {
	authOpts := swauth.AuthOpts{
		User: "test:tester",
		Key:  "testing",
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAuthSuccessfully(t, authOpts)

	providerClient, err := openstack.NewClient(th.Endpoint())
	th.AssertNoErr(t, err)

	swiftClient, err := swauth.NewObjectStorageV1(providerClient, authOpts)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, AuthResult.Token, swiftClient.TokenID)
}

func TestBadAuth(t *testing.T) {
	authOpts := swauth.AuthOpts{}
	_, err := authOpts.ToAuthOptsMap()
	if err == nil {
		t.Fatalf("Expected an error due to missing auth options")
	}
}
