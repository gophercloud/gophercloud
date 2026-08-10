package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// IdPTokenResponse is a client-credentials token response.
const IdPTokenResponse = `{
	"access_token": "test-idp-access-token-abc123",
	"token_type": "Bearer",
	"expires_in": 3600,
	"scope": "openid"
}`

// IdPTokenResponseIDToken contains an ID token instead of an access token.
const IdPTokenResponseIDToken = `{
	"id_token": "test-idp-id-token-xyz789",
	"token_type": "Bearer",
	"expires_in": 3600,
	"scope": "openid"
}`

// DiscoveryDocument is a minimal OIDC discovery document.
const DiscoveryDocument = `{
	"token_endpoint": "%s",
	"grant_types_supported": ["authorization_code", "client_credentials"]
}`

// DiscoveryDocumentNoGrantTypes omits grant_types_supported.
const DiscoveryDocumentNoGrantTypes = `{
	"issuer": "https://idp.example.com",
	"token_endpoint": "%s"
}`

// DiscoveryDocumentUnsupportedGrant omits client_credentials.
const DiscoveryDocumentUnsupportedGrant = `{
	"issuer": "https://idp.example.com",
	"token_endpoint": "https://idp.example.com/oauth2/token",
	"grant_types_supported": ["authorization_code", "refresh_token"]
}`

// DiscoveryDocumentNoTokenEndpoint omits token_endpoint.
const DiscoveryDocumentNoTokenEndpoint = `{
	"issuer": "https://idp.example.com",
	"grant_types_supported": ["client_credentials"]
}`

// FederationAuthResponse is an unscoped federation token response.
const FederationAuthResponse = `{
	"token": {
		"methods": ["mapped"],
		"user": {
			"id": "federation-user-id-001",
			"name": "federation-user",
			"domain": {
				"id": "Federated",
				"name": "Federated"
			},
			"OS-FEDERATION": {
				"identity_provider": {"id": "my-idp"},
				"protocol": {"id": "openid"},
				"groups": [
					{"id": "group-abc-001"}
				]
			}
		},
		"expires_at": "2017-06-03T02:19:49.000000Z",
		"is_domain": false
	}
}`

// UnscopedTokenID is the unscoped federation token ID.
const UnscopedTokenID = "unscoped-federation-token-aaa111"

// RescopedTokenResponse is a project-scoped token response.
const RescopedTokenResponse = `{
	"token": {
		"methods": ["token"],
		"project": {
			"domain": {
				"id": "default",
				"name": "Default"
			},
			"id": "project-id-001",
			"name": "my-project"
		},
		"catalog": [
			{
				"endpoints": [
					{
						"url": "http://127.0.0.1:8774/v2.1/project-id-001",
						"interface": "public",
						"region": "RegionOne",
						"region_id": "RegionOne",
						"id": "endpoint-compute-001"
					}
				],
				"type": "compute",
				"id": "service-compute-001",
				"name": "nova"
			}
		],
		"user": {
			"id": "federation-user-id-001",
			"name": "federation-user",
			"domain": {
				"id": "Federated",
				"name": "Federated"
			}
		},
		"roles": [
			{
				"id": "role-member-001",
				"name": "member"
			}
		],
		"expires_at": "2017-06-03T02:19:49.000000Z",
		"is_domain": false
	}
}`

// ScopedTokenID is the X-Subject-Token returned after rescoping.
const ScopedTokenID = "scoped-federation-token-bbb222"

// HandleIdPTokenSuccessfully handles a client-credentials grant.
func HandleIdPTokenSuccessfully(t *testing.T, fakeServer th.FakeServer, expectedBasicAuth string) {
	fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		th.TestHeader(t, r, "Accept", "application/json")

		if expectedBasicAuth != "" {
			th.TestHeader(t, r, "Authorization", expectedBasicAuth)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}
		if r.PostForm.Get("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", r.PostForm.Get("grant_type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, IdPTokenResponse)
	})
}

// HandleIdPTokenIDTokenSuccessfully returns an ID token.
func HandleIdPTokenIDTokenSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, IdPTokenResponseIDToken)
	})
}

// HandleFederationAuthSuccessfully handles a federation token exchange.
func HandleFederationAuthSuccessfully(t *testing.T, fakeServer th.FakeServer, expectedBearer string) {
	fakeServer.Mux.HandleFunc("/OS-FEDERATION/identity_providers/my-idp/protocols/openid/auth", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", "Bearer "+expectedBearer)

		if r.Header.Get("X-Auth-Token") != "" {
			t.Error("Expected X-Auth-Token header to be omitted")
		}

		w.Header().Set("X-Subject-Token", UnscopedTokenID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, FederationAuthResponse)
	})
}

// HandleRescopeSuccessfully handles token rescoping.
func HandleRescopeSuccessfully(t *testing.T, fakeServer th.FakeServer, expectedUnscopedToken string) {
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")

		if r.Header.Get("X-Auth-Token") != "" {
			t.Error("Expected X-Auth-Token header to be omitted")
		}

		w.Header().Set("X-Subject-Token", ScopedTokenID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, RescopedTokenResponse)
	})
}

// HandleDiscoverySuccessfully returns a discovery document.
func HandleDiscoverySuccessfully(t *testing.T, fakeServer th.FakeServer, tokenEndpointURL string) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, DiscoveryDocument, tokenEndpointURL)
	})
}

// HandleDiscoveryNoGrantTypesSuccessfully omits grant_types_supported.
func HandleDiscoveryNoGrantTypesSuccessfully(t *testing.T, fakeServer th.FakeServer, tokenEndpointURL string) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, DiscoveryDocumentNoGrantTypes, tokenEndpointURL)
	})
}

// HandleDiscoveryUnsupportedGrantType omits client_credentials.
func HandleDiscoveryUnsupportedGrantType(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, DiscoveryDocumentUnsupportedGrant)
	})
}

// HandleDiscoveryNoTokenEndpoint omits token_endpoint.
func HandleDiscoveryNoTokenEndpoint(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, DiscoveryDocumentNoTokenEndpoint)
	})
}

// HandleDiscoveryHTTPError returns an HTTP error.
func HandleDiscoveryHTTPError(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
	})
}

// HandleDiscoveryInvalidJSON returns invalid JSON.
func HandleDiscoveryInvalidJSON(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "not valid json{{{")
	})
}
