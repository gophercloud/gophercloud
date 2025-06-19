package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/applicationcredentials"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListApplicationCredentials(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListApplicationCredentialsSuccessfully(t, fakeServer)

	count := 0
	err := applicationcredentials.List(client.ServiceClient(fakeServer), userID, nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := applicationcredentials.ExtractApplicationCredentials(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedApplicationCredentialsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListApplicationCredentialsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListApplicationCredentialsSuccessfully(t, fakeServer)

	allPages, err := applicationcredentials.List(client.ServiceClient(fakeServer), userID, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := applicationcredentials.ExtractApplicationCredentials(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedApplicationCredentialsSlice, actual)
	th.AssertDeepEquals(t, ExpectedApplicationCredentialsSlice[0].Roles, []applicationcredentials.Role{{ID: "31f87923ae4a4d119aa0b85dcdbeed13", Name: "compute_viewer"}})
	th.AssertDeepEquals(t, ExpectedApplicationCredentialsSlice[1].Roles, []applicationcredentials.Role{{ID: "31f87923ae4a4d119aa0b85dcdbeed13", Name: "compute_viewer"}, {ID: "4494bc5bea1a4105ad7fbba6a7eb9ef4", Name: "network_viewer"}})
}

func TestGetApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetApplicationCredentialSuccessfully(t, fakeServer)

	actual, err := applicationcredentials.Get(context.TODO(), client.ServiceClient(fakeServer), userID, applicationCredentialID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ApplicationCredential, *actual)
}

func TestCreateApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateApplicationCredentialSuccessfully(t, fakeServer)

	createOpts := applicationcredentials.CreateOpts{
		Name:   "test",
		Secret: "mysecret",
		Roles: []applicationcredentials.Role{
			{ID: "31f87923ae4a4d119aa0b85dcdbeed13"},
		},
		AccessRules: []applicationcredentials.AccessRule{
			{
				Path:    "/v2.0/metrics",
				Method:  "GET",
				Service: "monitoring",
			},
		},
	}

	ApplicationCredentialResponse := ApplicationCredential
	ApplicationCredentialResponse.Secret = "mysecret"

	actual, err := applicationcredentials.Create(context.TODO(), client.ServiceClient(fakeServer), userID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ApplicationCredentialResponse, *actual)
}

func TestCreateNoSecretApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateNoSecretApplicationCredentialSuccessfully(t, fakeServer)

	createOpts := applicationcredentials.CreateOpts{
		Name: "test1",
		Roles: []applicationcredentials.Role{
			{ID: "31f87923ae4a4d119aa0b85dcdbeed13"},
		},
	}

	actual, err := applicationcredentials.Create(context.TODO(), client.ServiceClient(fakeServer), userID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ApplicationCredentialNoSecretResponse, *actual)
}

func TestCreateUnrestrictedApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateUnrestrictedApplicationCredentialSuccessfully(t, fakeServer)

	createOpts := applicationcredentials.CreateOpts{
		Name:         "test2",
		Unrestricted: true,
		Roles: []applicationcredentials.Role{
			{ID: "31f87923ae4a4d119aa0b85dcdbeed13"},
			{ID: "4494bc5bea1a4105ad7fbba6a7eb9ef4"},
		},
		ExpiresAt: &ApplationCredentialExpiresAt,
	}

	UnrestrictedApplicationCredentialResponse := UnrestrictedApplicationCredential
	UnrestrictedApplicationCredentialResponse.Secret = "generated_secret"

	actual, err := applicationcredentials.Create(context.TODO(), client.ServiceClient(fakeServer), userID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UnrestrictedApplicationCredentialResponse, *actual)
}

func TestDeleteApplicationCredential(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteApplicationCredentialSuccessfully(t, fakeServer)

	res := applicationcredentials.Delete(context.TODO(), client.ServiceClient(fakeServer), userID, applicationCredentialID)
	th.AssertNoErr(t, res.Err)
}

func TestListAccessRules(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAccessRulesSuccessfully(t, fakeServer)

	count := 0
	err := applicationcredentials.ListAccessRules(client.ServiceClient(fakeServer), userID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := applicationcredentials.ExtractAccessRules(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedAccessRulesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetAccessRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetAccessRuleSuccessfully(t, fakeServer)

	actual, err := applicationcredentials.GetAccessRule(context.TODO(), client.ServiceClient(fakeServer), userID, accessRuleID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AccessRule, *actual)
}

func TestDeleteAccessRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteAccessRuleSuccessfully(t, fakeServer)

	res := applicationcredentials.DeleteAccessRule(context.TODO(), client.ServiceClient(fakeServer), userID, accessRuleID)
	th.AssertNoErr(t, res.Err)
}
