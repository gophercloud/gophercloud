package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV3PasswordOptsToAuthScope(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	opts := auth.V3PasswordOpts{Username: "u", Password: "p", Scope: scope}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{"project": map[string]any{"id": &scope.ProjectID}}, scopeMap)
}

func TestV3TOTPOptsToAuthScope(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	opts := auth.V3TOTPOpts{Username: "u", Passcode: "123456", Scope: scope}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{"project": map[string]any{"id": &scope.ProjectID}}, scopeMap)
}

func TestV3TokenOptsToAuthScope(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	opts := auth.V3TokenOpts{Token: "t", Scope: scope}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{"project": map[string]any{"id": &scope.ProjectID}}, scopeMap)
}

func TestV3ApplicationCredentialOptsToAuthScope(t *testing.T) {
	opts := auth.V3ApplicationCredentialOpts{
		ApplicationCredentialID:     "id",
		ApplicationCredentialSecret: "secret",
		UserID:                      "user-id",
	}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, scopeMap == nil)
}

func TestV3MultifactorOptsToAuthScope(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id"}
	opts := auth.V3MultifactorOpts{Scope: scope}
	scopeMap, err := opts.ToAuthScope()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]any{"project": map[string]any{"id": &scope.ProjectID}}, scopeMap)
}
