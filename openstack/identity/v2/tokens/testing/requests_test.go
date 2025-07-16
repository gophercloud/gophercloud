package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func tokenPost(t *testing.T, options gophercloud.AuthOptions, requestJSON string) tokens.CreateResult {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleTokenPost(t, fakeServer, requestJSON)

	return tokens.Create(context.TODO(), client.ServiceClient(fakeServer), options)
}

func tokenPostErr(t *testing.T, options gophercloud.AuthOptions, expectedErr error) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleTokenPost(t, fakeServer, "")

	actualErr := tokens.Create(context.TODO(), client.ServiceClient(fakeServer), options).Err
	th.CheckDeepEquals(t, expectedErr, actualErr)
}

func TestCreateWithPassword(t *testing.T) {
	options := gophercloud.AuthOptions{
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
	options := gophercloud.AuthOptions{
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
	options := gophercloud.AuthOptions{
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

func TestRequireUsername(t *testing.T) {
	options := gophercloud.AuthOptions{
		Password: "thing",
	}

	tokenPostErr(t, options, gophercloud.ErrMissingInput{Argument: "Username"})
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
