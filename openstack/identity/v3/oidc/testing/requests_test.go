package testing

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oidc"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func stringPointer(s string) *string {
	return &s
}

func TestCreateUnscoped(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestCreateScoped(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestCreateScopedByProjectID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		Scope: tokens.Scope{
			ProjectID: "project-id-001",
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestCreateScopedByDomainID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		Scope: tokens.Scope{
			DomainID: "default",
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestCreateScopedByDomainName(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		Scope: tokens.Scope{
			DomainName: "Default",
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestCreateWithIDToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleIdPTokenIDTokenSuccessfully(t, fakeServer)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-id-token-xyz789")

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		AccessTokenType:      "id_token",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestCreateRejectsUnsupportedAccessTokenType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"opaque-token","token_type":"MAC"}`)
	})
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}
	result := oidc.Create(context.Background(), &client, &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
	})
	if result.Err == nil || !strings.Contains(result.Err.Error(), "unsupported token_type") {
		t.Fatalf("expected unsupported token_type error, got %v", result.Err)
	}
}

func TestCreateScopedBySystemScope(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		Scope: tokens.Scope{
			System: true,
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestCanReauth(t *testing.T) {
	opts := &oidc.AuthOptions{
		AllowReauth: true,
	}
	th.CheckEquals(t, true, opts.CanReauth())

	opts.AllowReauth = false
	th.CheckEquals(t, false, opts.CanReauth())
}

func TestValidationMissingInput(t *testing.T) {
	tests := []struct {
		argument string
		omit     func(*oidc.AuthOptions)
	}{
		{"IdentityProviderName", func(opts *oidc.AuthOptions) { opts.IdentityProviderName = "" }},
		{"Protocol", func(opts *oidc.AuthOptions) { opts.Protocol = "" }},
		{"ClientID", func(opts *oidc.AuthOptions) { opts.ClientID = "" }},
		{"AccessTokenEndpoint/DiscoveryEndpoint", func(opts *oidc.AuthOptions) { opts.AccessTokenEndpoint = "" }},
	}
	validators := map[string]func(*oidc.AuthOptions) error{
		"ToTokenV3CreateMap": func(opts *oidc.AuthOptions) error {
			_, err := opts.ToTokenV3CreateMap(nil)
			return err
		},
		"Create": func(opts *oidc.AuthOptions) error {
			return oidc.Create(context.Background(), nil, opts).Err
		},
	}
	for _, test := range tests {
		t.Run(test.argument, func(t *testing.T) {
			for name, validate := range validators {
				t.Run(name, func(t *testing.T) {
					opts := &oidc.AuthOptions{
						IdentityProviderName: "my-idp",
						Protocol:             "openid",
						ClientID:             "my-client-id",
						AccessTokenEndpoint:  "https://idp.example.com/oauth2/token",
					}
					test.omit(opts)
					err := validate(opts)
					var missing gophercloud.ErrMissingInput
					if !errors.As(err, &missing) {
						t.Fatalf("Expected ErrMissingInput, got %v", err)
					}
					th.CheckEquals(t, test.argument, missing.Argument)
				})
			}
		})
	}
}

func TestCreateUnscopedWithDiscovery(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	tokenEndpointURL := fakeServer.Endpoint() + "oauth2/token"
	HandleDiscoverySuccessfully(t, fakeServer, tokenEndpointURL)

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestCreateScopedWithDiscovery(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	tokenEndpointURL := fakeServer.Endpoint() + "oauth2/token"
	HandleDiscoverySuccessfully(t, fakeServer, tokenEndpointURL)

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")
	HandleRescopeSuccessfully(t, fakeServer, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ScopedTokenID, token)
}

func TestAccessTokenEndpointOverridesDiscovery(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
		DiscoveryEndpoint:    "http://should-not-be-called.example.com/.well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestDiscoveryWithNoGrantTypesField(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	tokenEndpointURL := fakeServer.Endpoint() + "oauth2/token"
	HandleDiscoveryNoGrantTypesSuccessfully(t, fakeServer, tokenEndpointURL)

	expectedBasicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-client-id:my-client-secret"))
	HandleIdPTokenSuccessfully(t, fakeServer, expectedBasicAuth)
	HandleFederationAuthSuccessfully(t, fakeServer, "test-idp-access-token-abc123")

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestDiscoveryUnsupportedGrantType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDiscoveryUnsupportedGrantType(t, fakeServer)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertErr(t, result.Err)
}

func TestDiscoveryNoTokenEndpoint(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDiscoveryNoTokenEndpoint(t, fakeServer)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertErr(t, result.Err)
}

func TestDiscoveryHTTPError(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDiscoveryHTTPError(t, fakeServer)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertErr(t, result.Err)
}

func TestDiscoveryInvalidJSON(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDiscoveryInvalidJSON(t, fakeServer)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		DiscoveryEndpoint:    fakeServer.Endpoint() + ".well-known/openid-configuration",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertErr(t, result.Err)
}

func TestValidationWrongOptionsType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &gophercloud.AuthOptions{
		Username: "user",
		Password: "pass",
	}

	result := oidc.Create(context.TODO(), &client, opts)
	th.AssertErr(t, result.Err)
}

func TestCreateRejectsNilServiceClient(t *testing.T) {
	result := oidc.Create(context.Background(), nil, &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "client-id",
		AccessTokenEndpoint:  "https://idp.example.com/token",
	})
	if result.Err == nil || !strings.Contains(result.Err.Error(), "ServiceClient") {
		t.Fatalf("expected nil ServiceClient error, got %v", result.Err)
	}
}

func TestToTokenV3ScopeMapProjectNameDomainName(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   stringPointer("my-project"),
			"domain": map[string]any{"name": stringPointer("Default")},
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapProjectNameDomainID(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainID:    "default",
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"name":   stringPointer("my-project"),
			"domain": map[string]any{"id": stringPointer("default")},
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapProjectID(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			ProjectID: "project-id-001",
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"project": map[string]any{
			"id": stringPointer("project-id-001"),
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapDomainID(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			DomainID: "default",
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"domain": map[string]any{
			"id": stringPointer("default"),
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapDomainName(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			DomainName: "Default",
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"domain": map[string]any{
			"name": stringPointer("Default"),
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapSystem(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			System: true,
		},
	}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	expected := map[string]any{
		"system": map[string]any{
			"all": true,
		},
	}
	th.CheckDeepEquals(t, expected, scope)
}

func TestToTokenV3ScopeMapUnscoped(t *testing.T) {
	opts := &oidc.AuthOptions{}

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)

	if scope != nil {
		t.Errorf("Expected nil scope for unscoped auth, got %v", scope)
	}
}

func TestToTokenV3ScopeMapProjectNameNoDomain(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			ProjectName: "my-project",
		},
	}

	_, err := opts.ToTokenV3ScopeMap()
	th.AssertErr(t, err)
}

func TestToTokenV3ScopeMapProjectNameAndID(t *testing.T) {
	opts := &oidc.AuthOptions{
		Scope: tokens.Scope{
			ProjectName: "my-project",
			ProjectID:   "project-id-001",
			DomainName:  "Default",
		},
	}

	_, err := opts.ToTokenV3ScopeMap()
	th.AssertErr(t, err)
}

func TestToTokenV3HeadersMap(t *testing.T) {
	opts := &oidc.AuthOptions{}

	headers, err := opts.ToTokenV3HeadersMap(nil)
	th.AssertNoErr(t, err)
	if headers != nil {
		t.Errorf("Expected nil headers, got %v", headers)
	}
}

func TestToTokenV3CreateMap(t *testing.T) {
	tests := []struct {
		name      string
		endpoint  string
		discovery string
	}{
		{name: "explicit endpoint", endpoint: "https://idp.example.com/oauth2/token"},
		{name: "discovery", discovery: "https://idp.example.com/.well-known/openid-configuration"},
		{name: "both endpoints", endpoint: "https://idp.example.com/oauth2/token", discovery: "https://idp.example.com/.well-known/openid-configuration"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := &oidc.AuthOptions{
				IdentityProviderName: "my-idp",
				Protocol:             "openid",
				ClientID:             "my-client-id",
				AccessTokenEndpoint:  test.endpoint,
				DiscoveryEndpoint:    test.discovery,
			}
			body, err := opts.ToTokenV3CreateMap(nil)
			th.AssertNoErr(t, err)
			if body != nil {
				t.Errorf("Expected nil body, got %v", body)
			}
		})
	}
}

var _ tokens.AuthOptionsBuilder = (*oidc.AuthOptions)(nil)
