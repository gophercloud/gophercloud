//go:build acceptance || identity || roles

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/roles"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/users"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRolesAddToUser(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2AdminClient()
	th.AssertNoErr(t, err)

	tenant, err := FindTenant(t, client)
	th.AssertNoErr(t, err)

	role, err := FindRole(t, client)
	th.AssertNoErr(t, err)

	user, err := CreateUser(t, client, tenant)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user)

	err = AddUserRole(t, client, tenant, user, role)
	th.AssertNoErr(t, err)
	defer DeleteUserRole(t, client, tenant, user, role)

	allPages, err := users.ListRoles(client, tenant.ID, user.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := users.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Roles of user %s:", user.Name)
	var found bool
	for _, r := range allRoles {
		tools.PrintResource(t, role)
		if r.Name == role.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestRolesList(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2AdminClient()
	th.AssertNoErr(t, err)

	allPages, err := roles.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRoles, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, r := range allRoles {
		tools.PrintResource(t, r)
		if r.Name == "admin" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
