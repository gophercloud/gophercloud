package testing

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oidc"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreateWithEncodedBasicCredentials(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	// OAuth client credentials are form-encoded before HTTP Basic encoding.
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("client%3Awith%2Breserved:secret%2Bwith%25reserved"))
	HandleIdPTokenSuccessfully(t, fakeServer, expected)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")

	client := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}
	result := oidc.Create(context.Background(), client, &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "client:with+reserved",
		ClientSecret:         "secret+with%reserved",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
	})
	th.AssertNoErr(t, result.Err)
	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}
