package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/roles"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListRoleResponse(t)

	count := 0

	err := roles.List(client.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := roles.ExtractRoles(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []roles.Role{
			{
				ID:          "123",
				Name:        "compute:admin",
				Description: "Nova Administrator",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestAddUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockAddUserRoleResponse(t)

	err := roles.AddUser(context.TODO(), client.ServiceClient(), "{tenant_id}", "{user_id}", "{role_id}").ExtractErr()

	th.AssertNoErr(t, err)
}

func TestDeleteUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteUserRoleResponse(t)

	err := roles.DeleteUser(context.TODO(), client.ServiceClient(), "{tenant_id}", "{user_id}", "{role_id}").ExtractErr()

	th.AssertNoErr(t, err)
}
