//go:build acceptance || identity || tokens

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	identityv2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTokenAuthenticate(t *testing.T) {
	identityv2.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2UnauthenticatedClient()
	th.AssertNoErr(t, err)

	authOptions, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)

	result, err := authOptions.Authenticate(context.TODO(), client.ProviderClient)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, result)

	for _, entry := range result.Catalog.Entries {
		tools.PrintResource(t, entry)
	}
}

func TestTokenValidate(t *testing.T) {
	identityv2.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2Client()
	th.AssertNoErr(t, err)

	authOptions, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)

	result, err := authOptions.Authenticate(context.TODO(), client.ProviderClient)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, result)

	getResult := tokens.Get(context.TODO(), client, result.TokenID)
	user, err := getResult.ExtractUser()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, user)
}
