package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/noauth"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestNoAuth(t *testing.T) {
	ao := gophercloud.AuthOptions{
		Username:   "user",
		TenantName: "test",
	}
	provider, err := noauth.UnAuthenticatedClient(ao)
	th.AssertNoErr(t, err)
	noauthClient, err := noauth.NewBlockStorageV2(provider, noauth.EndpointOpts{
		CinderEndpoint: "http://cinder:8776/v2",
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, noauthClient.Endpoint, naTestResult.Endpoint)
	th.AssertEquals(t, noauthClient.TokenID, naTestResult.TokenID)

	ao2 := gophercloud.AuthOptions{}
	provider2, err := noauth.UnAuthenticatedClient(ao2)
	th.AssertNoErr(t, err)
	noauthClient2, err := noauth.NewBlockStorageV2(provider2, noauth.EndpointOpts{
		CinderEndpoint: "http://cinder:8776/v2",
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, noauthClient2.Endpoint, naResult.Endpoint)
	th.AssertEquals(t, noauthClient2.TokenID, naResult.TokenID)
}
