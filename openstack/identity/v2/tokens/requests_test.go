package tokens

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func tokenPost(t *testing.T, options gophercloud.AuthOptions, requestJSON string) CreateResult {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenPost(t, requestJSON)
	return Create(client.ServiceClient(), options)
}

func tokenPostErr(t *testing.T, options gophercloud.AuthOptions, expectedErr error) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenPost(t, "")

	actualErr := Create(client.ServiceClient(), options).Err
	th.CheckDeepEquals(t, reflect.TypeOf(expectedErr), reflect.TypeOf(actualErr))
}

func TestCreateWithToken(t *testing.T) {
	options := gophercloud.AuthOptions{
		TokenID: "cbc36478b0bd8e67e89469c7749d4127",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "token": {
          "id": "cbc36478b0bd8e67e89469c7749d4127"
        }
      }
    }
  `))
}

func TestCreateWithPassword(t *testing.T) {
	options := gophercloud.AuthOptions{}
	options.Username = "me"
	options.Password = "swordfish"

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
	options := gophercloud.AuthOptions{}
	options.Username = "me"
	options.Password = "opensesame"
	options.TenantID = "fc394f2ab2df4114bde39905f800dc57"

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
	options := gophercloud.AuthOptions{}
	options.Username = "me"
	options.Password = "opensesame"
	options.TenantName = "demo"

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
	options := gophercloud.AuthOptions{}
	options.Password = "thing"

	expected := gophercloud.ErrMissingInput{}
	expected.Argument = "tokens.AuthOptions.Username/tokens.AuthOptions.TokenID"
	expected.Info = "You must provide either username/password or tenantID/token values."
	tokenPostErr(t, options, expected)
}

func TestRequirePassword(t *testing.T) {
	options := gophercloud.AuthOptions{}
	options.Username = "me"

	expected := gophercloud.ErrMissingInput{}
	expected.Argument = "tokens.AuthOptions.Password"
	tokenPostErr(t, options, expected)
}

func tokenGet(t *testing.T, tokenID string) GetResult {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenGet(t, tokenID)
	return Get(client.ServiceClient(), tokenID)
}

func TestGetWithToken(t *testing.T) {
	GetIsSuccessful(t, tokenGet(t, "db22caf43c934e6c829087c41ff8d8d6"))
}
