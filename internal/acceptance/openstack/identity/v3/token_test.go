//go:build acceptance || identity || tokens

package v3

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
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

	authData, err := ao.Auth.ToAuthBody()
	th.AssertNoErr(t, err)
	identity := map[string]any{"methods": slices.Sorted(maps.Keys(authData))}
	for method, body := range authData {
		identity[method] = body
	}
	authBody := map[string]any{"identity": identity}
	scope, err := ao.Auth.ToAuthScope()
	th.AssertNoErr(t, err)
	if scope != nil {
		authBody["scope"] = scope
	}
	authHeaders, err := ao.Auth.ToAuthHeaders()
	th.AssertNoErr(t, err)
	headers := make(map[string]string, len(authHeaders))
	for name, value := range authHeaders {
		headers[name] = fmt.Sprint(value)
	}
	request := gophercloud.RawRequest{Body: map[string]any{"auth": authBody}, Headers: headers}

	token, err := tokens.Create(context.TODO(), client, request).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	catalog, err := tokens.Get(context.TODO(), client, token.ID, nil).ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, catalog)

	user, err := tokens.Get(context.TODO(), client, token.ID, nil).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokens.Get(context.TODO(), client, token.ID, nil).ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokens.Get(context.TODO(), client, token.ID, nil).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)
}
