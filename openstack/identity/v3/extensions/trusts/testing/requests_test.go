package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateUserIDPasswordTrustID(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	ao := trusts.AuthOptsExt{
		TrustID: "de0945a",
		AuthOptionsBuilder: &tokens.AuthOptions{
			UserID:   "me",
			Password: "squirrel!",
		},
	}
	HandleCreateTokenWithTrustID(t, ao, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": { "id": "me", "password": "squirrel!" }
					}
				},
        "scope": {
            "OS-TRUST:trust": {
                "id": "de0945a"
            }
        }
			}
		}
	`)

	var actual struct {
		tokens.Token
		trusts.TokenExt
	}
	err := tokens.Create(client.ServiceClient(), ao).ExtractInto(&actual)
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}
	expected := struct {
		tokens.Token
		trusts.TokenExt
	}{
		tokens.Token{
			ExpiresAt: time.Date(2013, 02, 27, 18, 30, 59, 999999000, time.UTC),
		},
		trusts.TokenExt{
			Trust: trusts.Trust{
				ID:                 "fe0aef",
				Impersonation:      false,
				RedelegatedTrustID: "3ba234",
				RedelegationCount:  2,
			},
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateTrust(t)

	expiresAt := time.Date(2019, 12, 1, 14, 0, 0, 999999999, time.UTC)
	result, err := trusts.Create(client.ServiceClient(), trusts.CreateOpts{
		ExpiresAt:         &expiresAt,
		AllowRedelegation: true,
		ProjectID:         "9b71012f5a4a4aef9193f1995fe159b2",
		Roles: []trusts.Role{
			{
				Name: "member",
			},
		},
		TrusteeUserID: "ecb37e88cc86431c99d0332208cb6fbf",
		TrustorUserID: "959ed913a32c4ec88c041c98e61cbbc3",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "3422b7c113894f5d90665e1a79655e23", result.ID)
	th.AssertEquals(t, false, result.Impersonation)
	th.AssertEquals(t, 10, result.RedelegationCount)
}

func TestDeleteTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteTrust(t)

	res := trusts.Delete(client.ServiceClient(), "3422b7c113894f5d90665e1a79655e23")
	th.AssertNoErr(t, res.Err)
}

func TestGetTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTrustSuccessfully(t)

	res := trusts.Get(client.ServiceClient(), "987fe8")
	th.AssertNoErr(t, res.Err)
}

func TestListTrusts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustsSuccessfully(t)

	count := 0
	err := trusts.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := trusts.ExtractTrusts(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedTrustsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListTrustsAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustsSuccessfully(t)

	allPages, err := trusts.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := trusts.ExtractTrusts(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTrustsSlice, actual)
}

func TestListTrustsFiltered(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustsSuccessfully(t)
	trustsListOpts := trusts.ListOpts{
		TrustorUserID: "86c0d5",
	}
	allPages, err := trusts.List(client.ServiceClient(), trustsListOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := trusts.ExtractTrusts(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTrustsSlice, actual)
}
