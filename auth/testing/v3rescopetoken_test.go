package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV3TokenOptsRescopeBodyAndScope(t *testing.T) {
	opts := auth.V3TokenOpts{Token: "the-token", Scope: &auth.Scope{ProjectID: "project-id"}}
	body, _, err := v3AuthRequest(t, opts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"token"},
				"token":   map[string]any{"id": "the-token"},
			},
			"scope": map[string]any{"project": map[string]any{"id": "project-id"}},
		},
	}, body)
}

func TestV3TokenOptsUnscoped(t *testing.T) {
	body, _, err := v3AuthRequest(t, auth.V3TokenOpts{Token: "the-token"})
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"token"},
				"token":   map[string]any{"id": "the-token"},
			},
		},
	}, body)
}

func TestV3TokenOptsRescopeMissingToken(t *testing.T) {
	_, err := (auth.V3TokenOpts{Scope: &auth.Scope{ProjectID: "project-id"}}).ToAuthBody()
	th.AssertErr(t, err)
	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Token", missing.Argument)
}

func TestV3TokenOptsRescopeCannotReauthByDefault(t *testing.T) {
	opts := auth.V3TokenOpts{Token: "t", Scope: &auth.Scope{ProjectID: "p"}}
	th.AssertEquals(t, false, opts.CanReauth())
}
