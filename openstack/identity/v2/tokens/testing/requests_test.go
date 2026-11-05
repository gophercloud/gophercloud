package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func tokenPost(t *testing.T, options auth.AuthOptionsBuilderV2, requestJSON string) tokens.CreateResult {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleTokenPost(t, fakeServer, requestJSON)

	data, err := options.ToAuthBody()
	th.AssertNoErr(t, err)
	var body map[string]any
	for _, methodBody := range data {
		body = methodBody
	}
	return tokens.Create(context.TODO(), client.ServiceClient(fakeServer), gophercloud.RawRequest{Body: map[string]any{"auth": body}})
}

func TestCreateWithPassword(t *testing.T) {
	options := auth.V2PasswordOpts{
		Username: "me",
		Password: "swordfish",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "passwordCredentials": {
          "username": "me",
          "password": "swordfish"
        }
      }
    }
  `))
}

func TestCreateTokenWithTenantID(t *testing.T) {
	options := auth.V2PasswordOpts{
		Username: "me",
		Password: "opensesame",
		TenantID: "fc394f2ab2df4114bde39905f800dc57",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "tenantId": "fc394f2ab2df4114bde39905f800dc57",
        "passwordCredentials": {
          "username": "me",
          "password": "opensesame"
        }
      }
    }
  `))
}

func TestCreateTokenWithTenantName(t *testing.T) {
	options := auth.V2PasswordOpts{
		Username:   "me",
		Password:   "opensesame",
		TenantName: "demo",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "tenantName": "demo",
        "passwordCredentials": {
          "username": "me",
          "password": "opensesame"
        }
      }
    }
  `))
}

func tokenGet(t *testing.T, tokenId string) tokens.GetResult {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleTokenGet(t, fakeServer, tokenId)
	return tokens.Get(context.TODO(), client.ServiceClient(fakeServer), tokenId)
}

func TestGetWithToken(t *testing.T) {
	GetIsSuccessful(t, tokenGet(t, "db22caf43c934e6c829087c41ff8d8d6"))
}
