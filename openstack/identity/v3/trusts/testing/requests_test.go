package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/trusts"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateTrust(t)

	expiresAt := time.Date(2019, 12, 1, 14, 0, 0, 0, time.UTC)
	result, err := trusts.Create(context.TODO(), client.ServiceClient(), trusts.CreateOpts{
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

	th.AssertDeepEquals(t, CreatedTrust, *result)
}

func TestCreateTrustNoExpire(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateTrustNoExpire(t)

	result, err := trusts.Create(context.TODO(), client.ServiceClient(), trusts.CreateOpts{
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

	th.AssertDeepEquals(t, CreatedTrustNoExpire, *result)
}

func TestDeleteTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteTrust(t)

	res := trusts.Delete(context.TODO(), client.ServiceClient(), "3422b7c113894f5d90665e1a79655e23")
	th.AssertNoErr(t, res.Err)
}

func TestGetTrust(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTrustSuccessfully(t)

	res := trusts.Get(context.TODO(), client.ServiceClient(), "987fe8")
	th.AssertNoErr(t, res.Err)
}

func TestListTrusts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustsSuccessfully(t)

	count := 0
	err := trusts.List(client.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

	allPages, err := trusts.List(client.ServiceClient(), nil).AllPages(context.TODO())
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
	allPages, err := trusts.List(client.ServiceClient(), trustsListOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := trusts.ExtractTrusts(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTrustsSlice, actual)
}

func TestListTrustRoles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustRolesSuccessfully(t)

	count := 0
	err := trusts.ListRoles(client.ServiceClient(), "987fe8").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := trusts.ExtractRoles(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedTrustRolesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListTrustRolesAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTrustRolesSuccessfully(t)

	allPages, err := trusts.ListRoles(client.ServiceClient(), "987fe8").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := trusts.ExtractRoles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTrustRolesSlice, actual)
}

func TestGetTrustRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTrustRoleSuccessfully(t)

	role, err := trusts.GetRole(context.TODO(), client.ServiceClient(), "987fe8", "c1648e").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, FirstRole, *role)
}

func TestCheckTrustRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCheckTrustRoleSuccessfully(t)

	err := trusts.CheckRole(context.TODO(), client.ServiceClient(), "987fe8", "c1648e").ExtractErr()
	th.AssertNoErr(t, err)
}
