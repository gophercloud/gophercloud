//go:build acceptance || identity || tokens

package v3

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTokensCreateRawRequest(t *testing.T) {
	username, password := os.Getenv("OS_USERNAME"), os.Getenv("OS_PASSWORD")
	if username == "" || password == "" {
		t.Skip("password authentication requires OS_USERNAME and OS_PASSWORD")
	}

	client, err := clients.NewIdentityV3UnauthenticatedClient()
	th.AssertNoErr(t, err)

	userDomain := map[string]any{"id": "default"}
	if id := os.Getenv("OS_USER_DOMAIN_ID"); id != "" {
		userDomain = map[string]any{"id": id}
	} else if name := os.Getenv("OS_USER_DOMAIN_NAME"); name != "" {
		userDomain = map[string]any{"name": name}
	} else if id := os.Getenv("OS_DOMAIN_ID"); id != "" {
		userDomain = map[string]any{"id": id}
	} else if name := os.Getenv("OS_DOMAIN_NAME"); name != "" {
		userDomain = map[string]any{"name": name}
	}

	request := gophercloud.RawRequest{
		Body: map[string]any{
			"auth": map[string]any{
				"identity": map[string]any{
					"methods": []string{"password"},
					"password": map[string]any{
						"user": map[string]any{
							"name":     username,
							"password": password,
							"domain":   userDomain,
						},
					},
				},
			},
		},
		Headers: map[string]string{"Accept": "application/json"},
	}

	tokenID, err := tokens.Create(context.TODO(), client, request).ExtractTokenID()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, tokenID != "")
}
