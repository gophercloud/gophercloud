package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
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

func TestV3AuthMethodsPropagateScopeErrors(t *testing.T) {
	scope := &auth.Scope{ProjectID: "project-id", DomainID: "domain-id"}
	tests := []struct {
		name string
		opts auth.AuthOptionsBuilderV3
	}{
		{name: "password", opts: auth.V3PasswordOpts{Scope: scope}},
		{name: "TOTP", opts: auth.V3TOTPOpts{Scope: scope}},
		{name: "token", opts: auth.V3TokenOpts{Scope: scope}},
		{name: "rescope token", opts: auth.V3RescopeTokenOpts{Scope: scope}},
		{name: "multifactor", opts: auth.V3MultifactorOpts{Scope: scope}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthScope()
			th.AssertErr(t, err)
			_, ok := err.(gophercloud.ErrProjectScopeAndDomainScope)
			th.AssertEquals(t, true, ok)
		})
	}
}
