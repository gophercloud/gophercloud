//go:build acceptance
// +build acceptance

package v3

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack"
	"github.com/bizflycloud/gophercloud/openstack/identity/v3/tokens"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestTokensGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		Password:   ao.Password,
		DomainName: "default",
	}

	token, err := tokens.Create(client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	catalog, err := tokens.Get(client, token.ID).ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, catalog)

	user, err := tokens.Get(client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokens.Get(client, token.ID).ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokens.Get(client, token.ID).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)
}
