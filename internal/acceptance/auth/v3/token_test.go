//go:build acceptance || identity || tokens

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTokensGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	result, err := ao.Authenticate(context.TODO(), client.ProviderClient)
	th.AssertNoErr(t, err)

	token := result
	tools.PrintResource(t, token)

	catalog, err := tokens.Get(context.TODO(), client, token.TokenID, nil).ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, catalog)

	user, err := tokens.Get(context.TODO(), client, token.TokenID, nil).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokens.Get(context.TODO(), client, token.TokenID, nil).ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokens.Get(context.TODO(), client, token.TokenID, nil).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)
}
