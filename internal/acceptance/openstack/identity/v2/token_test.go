//go:build acceptance || identity || tokens

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTokenAuthenticate(t *testing.T) {
	RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2UnauthenticatedClient()
	th.AssertNoErr(t, err)

	authOptions, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	result, err := authOptions.Authenticate(context.TODO(), &client.HTTPClient)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)

	for _, entry := range result.Catalog.Entries {
		tools.PrintResource(t, entry)
	}
}

func TestTokenValidate(t *testing.T) {
	RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2Client()
	th.AssertNoErr(t, err)

	getResult := tokens.Get(context.TODO(), client, client.TokenID)
	user, err := getResult.ExtractUser()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, user)
}
