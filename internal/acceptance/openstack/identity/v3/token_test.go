//go:build acceptance || identity || tokens

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTokensGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	tokenResult := tokens.Get(context.TODO(), client, client.TokenID, nil)

	catalog, err := tokenResult.ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, catalog)

	user, err := tokenResult.ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokenResult.ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokenResult.ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)
}
