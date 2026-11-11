//go:build acceptance || identity || tokens

package v2

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// tokenCreateOpts supplies an already-assembled request to the Identity API.
type tokenCreateOpts struct {
	Body    map[string]any
	Headers map[string]string
}

func (opts tokenCreateOpts) ToTokenCreateMap() (map[string]any, error) {
	return opts.Body, nil
}

func (opts tokenCreateOpts) ToTokenCreateHeaders(method, url string) (map[string]string, error) {
	return opts.Headers, nil
}

func TestTokensCreateBuilder(t *testing.T) {
	RequireIdentityV2(t)

	username, password := os.Getenv("OS_USERNAME"), os.Getenv("OS_PASSWORD")
	if username == "" || password == "" {
		t.Skip("password authentication requires OS_USERNAME and OS_PASSWORD")
	}

	client, err := clients.NewIdentityV2UnauthenticatedClient()
	th.AssertNoErr(t, err)

	request := tokenCreateOpts{
		Body: map[string]any{
			"auth": map[string]any{
				"passwordCredentials": map[string]any{
					"username": username,
					"password": password,
				},
			},
		},
		Headers: map[string]string{"Accept": "application/json"},
	}

	tokenID, err := tokens.Create(context.TODO(), client, request).ExtractTokenID()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, tokenID != "")
}
