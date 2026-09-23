//go:build acceptance || identity || tokens

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
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

	authOptions, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)

	authData, err := authOptions.Auth.ToAuthBody()
	th.AssertNoErr(t, err)
	if len(authData) != 1 {
		t.Fatalf("expected one v2 auth method, got %d", len(authData))
	}
	var authBody map[string]any
	for _, body := range authData {
		authBody = body
	}
	request := gophercloud.RawRequest{Body: map[string]any{"auth": authBody}}
	result := tokens.Create(context.TODO(), client, request)
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
	RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2Client()
	th.AssertNoErr(t, err)

	authOptions, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)

	authData, err := authOptions.Auth.ToAuthBody()
	th.AssertNoErr(t, err)
	if len(authData) != 1 {
		t.Fatalf("expected one v2 auth method, got %d", len(authData))
	}
	var authBody map[string]any
	for _, body := range authData {
		authBody = body
	}
	request := gophercloud.RawRequest{Body: map[string]any{"auth": authBody}}
	result := tokens.Create(context.TODO(), client, request)
	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, token)

	getResult := tokens.Get(context.TODO(), client, token.ID)
	user, err := getResult.ExtractUser()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, user)
}
