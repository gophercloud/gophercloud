package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/swauth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuth(t *testing.T) {
	authOpts := swauth.AuthOpts{
		User: "test:tester",
		Key:  "testing",
	}

	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAuthSuccessfully(t, fakeServer, authOpts)

	providerClient, err := openstack.NewClient(fakeServer.Endpoint())
	th.AssertNoErr(t, err)

	swiftClient, err := swauth.NewObjectStorageV1(context.TODO(), providerClient, authOpts)
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
