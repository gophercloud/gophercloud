package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV3RescopeTokenOptsToAuthBody(t *testing.T) {
	opts := auth.V3RescopeTokenOpts{
		Token: "the-token",
		Scope: &auth.Scope{ProjectID: "project-id"},
	}

	authData, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)

	expected := map[string]map[string]any{
		"token": {"id": "the-token"},
	}
	th.AssertDeepEquals(t, expected, authData)
}

func TestV3RescopeTokenOptsMissingScope(t *testing.T) {
	opts := auth.V3RescopeTokenOpts{Token: "the-token"}

	_, err := opts.ToAuthBody()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeEmpty) // nolint: errorlint
	th.AssertEquals(t, true, ok)
}

func TestV3RescopeTokenOptsCanReauth(t *testing.T) {
	opts := auth.V3RescopeTokenOpts{Token: "t", Scope: &auth.Scope{ProjectID: "p"}}
	th.AssertEquals(t, false, opts.CanReauth())
}

func TestV3RescopeTokenOptsToAuthScope(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	opts := auth.V3RescopeTokenOpts{Token: "t", Scope: scope}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{"project": map[string]any{"id": &scope.ProjectID}}, scopeMap)
}
