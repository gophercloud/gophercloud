//go:build acceptance || identity
// +build acceptance identity

package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack"
	"github.com/bizflycloud/gophercloud/openstack/identity/v2/tokens"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestTokenAuthenticate(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2UnauthenticatedClient()
	th.AssertNoErr(t, err)

	authOptions, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	result := tokens.Create(client, authOptions)
	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, token)

	catalog, err := result.ExtractServiceCatalog()
	th.AssertNoErr(t, err)

	for _, entry := range catalog.Entries {
		tools.PrintResource(t, entry)
	}
}

func TestTokenValidate(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2Client()
	th.AssertNoErr(t, err)

	authOptions, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	result := tokens.Create(client, authOptions)
	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, token)

	getResult := tokens.Get(client, token.ID)
	user, err := getResult.ExtractUser()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, user)
}
