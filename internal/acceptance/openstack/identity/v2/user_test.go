//go:build acceptance || identity || users

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/users"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestUsersList(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2AdminClient()
	th.AssertNoErr(t, err)

	allPages, err := users.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allUsers, err := users.ExtractUsers(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, user := range allUsers {
		tools.PrintResource(t, user)

		if user.Name == "admin" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestUsersCreateUpdateDelete(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2AdminClient()
	th.AssertNoErr(t, err)

	tenant, err := FindTenant(t, client)
	th.AssertNoErr(t, err)

	user, err := CreateUser(t, client, tenant)
	th.AssertNoErr(t, err)
	defer DeleteUser(t, client, user)

	tools.PrintResource(t, user)

	newUser, err := UpdateUser(t, client, user)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newUser)
}
