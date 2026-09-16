package testing

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const (
	oidcAccessToken      = "idp-access-token"
	oidcTokenResponse    = `{"access_token": "idp-access-token", "token_type": "Bearer", "expires_in": 3600}`
	oidcUnscopedToken    = "unscoped-federation-token"
	oidcScopedToken      = "scoped-federation-token"
	oidcKeystoneResponse = `{"token": {"expires_at": "2099-01-01T00:00:00Z", "user": {"id": "user-id", "name": "federated-user"}}}`
	oidcFederationPath   = "/v3/OS-FEDERATION/identity_providers/my-idp/protocols/openid/auth"
)

func oidcBasicAuth(id, secret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(id+":"+secret))
}

// oidcOpts returns options for a fake server that is both the identity
// provider and Keystone.
func oidcOpts(fakeServer th.FakeServer) *auth.V3OIDCClientCredentialsOpts {
	return &auth.V3OIDCClientCredentialsOpts{
		AuthURL:              fakeServer.Endpoint(),
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  fakeServer.Endpoint() + "oauth2/token",
	}
}

func handleOIDCToken(t *testing.T, fakeServer th.FakeServer, status int, response string) {
	fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		th.CheckNoErr(t, r.ParseForm())
		th.CheckEquals(t, "client_credentials", r.PostForm.Get("grant_type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, response)
	})
}

func handleOIDCFederation(t *testing.T, fakeServer th.FakeServer, accessToken string) {
	fakeServer.Mux.HandleFunc(oidcFederationPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Authorization", "Bearer "+accessToken)
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		w.Header().Set("X-Subject-Token", oidcUnscopedToken)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, oidcKeystoneResponse)
	})
}

func TestV3OIDCClientCredentialsAuthenticate(t *testing.T) {
	tests := []struct {
		name  string
		scope *auth.Scope
		body  string
	}{
		{name: "unscoped"},
		{name: "project name", scope: &auth.Scope{ProjectName: "my-project", ProjectDomainName: "Default"}, body: `{"project": {"name": "my-project", "domain": {"name": "Default"}}}`},
		{name: "project ID", scope: &auth.Scope{ProjectID: "project-id"}, body: `{"project": {"id": "project-id"}}`},
		{name: "domain ID", scope: &auth.Scope{DomainID: "domain-id"}, body: `{"domain": {"id": "domain-id"}}`},
		{name: "domain name", scope: &auth.Scope{DomainName: "Default"}, body: `{"domain": {"name": "Default"}}`},
		{name: "system", scope: &auth.Scope{System: true}, body: `{"system": {"all": true}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			handleOIDCToken(t, fakeServer, http.StatusOK, oidcTokenResponse)
			handleOIDCFederation(t, fakeServer, oidcAccessToken)
			fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)
				th.TestHeaderUnset(t, r, "X-Auth-Token")
				th.TestJSONRequest(t, r, fmt.Sprintf(`{"auth": {"identity": {"methods": ["token"], "token": {"id": %q}}, "scope": %s}}`, oidcUnscopedToken, tt.body))
				w.Header().Set("X-Subject-Token", oidcScopedToken)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, oidcKeystoneResponse)
			})

			opts := oidcOpts(fakeServer)
			opts.Scope = tt.scope
			opts.AllowReauth = true
			provider := &gophercloud.ProviderClient{TokenID: "provider-token"}

			result, err := opts.Authenticate(context.Background(), provider)
			th.AssertNoErr(t, err)
			if tt.scope == nil {
				th.AssertEquals(t, oidcUnscopedToken, result.TokenID)
			} else {
				th.AssertEquals(t, oidcScopedToken, result.TokenID)
			}
			th.AssertEquals(t, true, result.CanReauth)
			th.AssertEquals(t, "provider-token", provider.Token())
		})
	}
}

func TestV3OIDCClientCredentialsTokenRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(*auth.V3OIDCClientCredentialsOpts)
		authorization string
		clientID      string
		scope         string
		response      string
		accessToken   string
	}{
		{
			name:          "defaults",
			modify:        func(*auth.V3OIDCClientCredentialsOpts) {},
			authorization: oidcBasicAuth("my-client-id", "my-client-secret"),
			scope:         "openid",
		},
		{
			// Like keystoneauth, the credentials are sent without form encoding.
			name: "reserved characters in credentials",
			modify: func(opts *auth.V3OIDCClientCredentialsOpts) {
				opts.ClientID, opts.ClientSecret = "client+with/reserved", "secret+with%reserved/="
			},
			authorization: oidcBasicAuth("client+with/reserved", "secret+with%reserved/="),
			scope:         "openid",
		},
		{
			name:     "no client secret",
			modify:   func(opts *auth.V3OIDCClientCredentialsOpts) { opts.ClientSecret = "" },
			clientID: "my-client-id",
			scope:    "openid",
		},
		{
			name:          "custom scope",
			modify:        func(opts *auth.V3OIDCClientCredentialsOpts) { opts.OIDCScope = "openid profile" },
			authorization: oidcBasicAuth("my-client-id", "my-client-secret"),
			scope:         "openid profile",
		},
		{
			name:          "lowercase token type",
			modify:        func(*auth.V3OIDCClientCredentialsOpts) {},
			authorization: oidcBasicAuth("my-client-id", "my-client-secret"),
			scope:         "openid",
			response:      `{"access_token": "idp-access-token", "token_type": "bearer"}`,
			accessToken:   oidcAccessToken,
		},
		{
			name:          "ID token",
			modify:        func(opts *auth.V3OIDCClientCredentialsOpts) { opts.AccessTokenType = "id_token" },
			authorization: oidcBasicAuth("my-client-id", "my-client-secret"),
			scope:         "openid",
			response:      `{"access_token": "idp-access-token", "id_token": "idp-id-token"}`,
			accessToken:   "idp-id-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			if tt.response == "" {
				tt.response, tt.accessToken = oidcTokenResponse, oidcAccessToken
			}
			fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
				th.CheckEquals(t, tt.authorization, r.Header.Get("Authorization"))
				th.CheckNoErr(t, r.ParseForm())
				th.CheckEquals(t, tt.clientID, r.PostForm.Get("client_id"))
				th.CheckEquals(t, tt.scope, r.PostForm.Get("scope"))
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.response)
			})
			handleOIDCFederation(t, fakeServer, tt.accessToken)

			opts := oidcOpts(fakeServer)
			tt.modify(opts)
			result, err := opts.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, oidcUnscopedToken, result.TokenID)
			th.AssertEquals(t, false, result.CanReauth)
		})
	}
}

func TestV3OIDCClientCredentialsRejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name             string
		tokenStatus      int
		tokenResponse    string
		federationStatus int
		err              string
	}{
		{name: "invalid JSON", tokenResponse: "not JSON", err: "failed to parse IdP token response"},
		{name: "missing access token", tokenResponse: `{"token_type": "Bearer"}`, err: `missing "access_token" field`},
		{name: "unsupported token type", tokenResponse: `{"access_token": "token", "token_type": "MAC"}`, err: `unsupported token_type "MAC"`},
		{name: "missing token type", tokenResponse: `{"access_token": "token"}`, err: "unsupported token_type"},
		{name: "IdP error", tokenStatus: http.StatusBadRequest, tokenResponse: `{"error": "invalid_client"}`, err: `IdP returned HTTP 400: {"error": "invalid_client"}`},
		{name: "long IdP error", tokenStatus: http.StatusBadRequest, tokenResponse: strings.Repeat("x", 600), err: "HTTP 400: " + strings.Repeat("x", 512) + "..."},
		{name: "oversized response", tokenResponse: strings.Repeat(" ", 1<<20) + oidcTokenResponse, err: "IdP response exceeds 1048576 bytes"},
		{name: "rejected by Keystone", tokenResponse: oidcTokenResponse, federationStatus: http.StatusUnauthorized},
		{name: "missing subject token", tokenResponse: oidcTokenResponse, federationStatus: http.StatusCreated, err: "missing X-Subject-Token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			if tt.tokenStatus == 0 {
				tt.tokenStatus = http.StatusOK
			}
			handleOIDCToken(t, fakeServer, tt.tokenStatus, tt.tokenResponse)
			fakeServer.Mux.HandleFunc(oidcFederationPath, func(w http.ResponseWriter, r *http.Request) {
				if tt.federationStatus == 0 {
					t.Error("unexpected request to Keystone")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.federationStatus)
				fmt.Fprint(w, oidcKeystoneResponse)
			})

			_, err := oidcOpts(fakeServer).Authenticate(context.Background(), nil)
			th.AssertErr(t, err)
			if tt.err == "" {
				th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, tt.federationStatus))
				return
			}
			if !strings.Contains(err.Error(), tt.err) {
				t.Fatalf("expected error containing %q, got %.200q", tt.err, err)
			}
			th.AssertEquals(t, true, len(err.Error()) < 1024)
		})
	}
}

func TestV3OIDCClientCredentialsDiscovery(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		document string
		modify   func(*auth.V3OIDCClientCredentialsOpts)
		err      string
	}{
		{name: "token endpoint", document: `{"token_endpoint": %q, "grant_types_supported": ["authorization_code", "client_credentials"]}`},
		{name: "no grant types", document: `{"token_endpoint": %q}`},
		{
			name:     "access token endpoint takes precedence",
			status:   http.StatusInternalServerError,
			document: `{}`,
			modify: func(opts *auth.V3OIDCClientCredentialsOpts) {
				opts.AccessTokenEndpoint = strings.TrimSuffix(opts.DiscoveryEndpoint, "discovery") + "oauth2/token"
			},
		},
		{name: "unsupported grant type", document: `{"token_endpoint": %q, "grant_types_supported": ["authorization_code"]}`, err: "does not support client_credentials"},
		{name: "no token endpoint", document: `{"issuer": "https://idp.example.com"}`, err: "does not contain a valid token_endpoint"},
		{name: "HTTP error", status: http.StatusInternalServerError, document: `{}`, err: "discovery endpoint returned HTTP 500"},
		{name: "invalid JSON", document: "not JSON", err: "discovery document is not valid JSON"},
		{name: "oversized document", document: strings.Repeat(" ", 1<<20) + `{"token_endpoint": %q}`, err: "discovery document exceeds 1048576 bytes"},
		{name: "plaintext token endpoint", document: `{"token_endpoint": "http://idp.example.com/oauth2/token"}`, err: `discovery document token_endpoint "http://idp.example.com/oauth2/token" does not use https`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			if tt.status == 0 {
				tt.status = http.StatusOK
			}
			if strings.Contains(tt.document, "%q") {
				tt.document = fmt.Sprintf(tt.document, fakeServer.Endpoint()+"oauth2/token")
			}
			fakeServer.Mux.HandleFunc("/discovery", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)
				th.TestHeader(t, r, "Accept", "application/json")
				th.TestHeaderUnset(t, r, "Authorization")
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.document)
			})
			handleOIDCToken(t, fakeServer, http.StatusOK, oidcTokenResponse)
			handleOIDCFederation(t, fakeServer, oidcAccessToken)

			opts := oidcOpts(fakeServer)
			opts.AccessTokenEndpoint = ""
			opts.DiscoveryEndpoint = fakeServer.Endpoint() + "discovery"
			if tt.modify != nil {
				tt.modify(opts)
			}
			result, err := opts.Authenticate(context.Background(), nil)
			if tt.err == "" {
				th.AssertNoErr(t, err)
				th.AssertEquals(t, oidcUnscopedToken, result.TokenID)
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.err) {
				t.Fatalf("expected error containing %q, got %.200q", tt.err, err)
			}
		})
	}
}

func TestV3OIDCClientCredentialsRequiresInput(t *testing.T) {
	tests := []struct {
		argument string
		omit     func(*auth.V3OIDCClientCredentialsOpts)
	}{
		{"AuthURL", func(opts *auth.V3OIDCClientCredentialsOpts) { opts.AuthURL = "" }},
		{"IdentityProviderName", func(opts *auth.V3OIDCClientCredentialsOpts) { opts.IdentityProviderName = "" }},
		{"Protocol", func(opts *auth.V3OIDCClientCredentialsOpts) { opts.Protocol = "" }},
		{"ClientID", func(opts *auth.V3OIDCClientCredentialsOpts) { opts.ClientID = "" }},
		{"AccessTokenEndpoint/DiscoveryEndpoint", func(opts *auth.V3OIDCClientCredentialsOpts) { opts.AccessTokenEndpoint = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.argument, func(t *testing.T) {
			opts := &auth.V3OIDCClientCredentialsOpts{
				AuthURL:              "https://identity.example.com/v3",
				IdentityProviderName: "my-idp",
				Protocol:             "openid",
				ClientID:             "my-client-id",
				AccessTokenEndpoint:  "https://idp.example.com/oauth2/token",
			}
			tt.omit(opts)
			_, err := opts.Authenticate(context.Background(), nil)
			var missing gophercloud.ErrMissingInput
			th.AssertEquals(t, true, errors.As(err, &missing))
			th.AssertEquals(t, tt.argument, missing.Argument)
		})
	}

	var opts *auth.V3OIDCClientCredentialsOpts
	_, err := opts.Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
}

// The identity provider endpoints carry the client secret, or say where to
// send it, so none of them may be plaintext off this machine, and a remote
// discovery document may not name a plaintext port on it.
func TestV3OIDCClientCredentialsRejectsPlaintextEndpoints(t *testing.T) {
	var plaintextRequests atomic.Int32
	local := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plaintextRequests.Add(1)
	}))
	defer local.Close()
	remote := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"token_endpoint": %q}`, local.URL+"/token")
	}))
	defer remote.Close()
	// The test server's certificate also covers example.com, which stands in
	// for a host that is not loopback.
	transport := remote.Client().Transport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	transport.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, remote.Listener.Addr().String())
	}
	provider := &gophercloud.ProviderClient{HTTPClient: http.Client{Transport: transport}}

	tests := []struct {
		name     string
		opts     auth.V3OIDCClientCredentialsOpts
		rejected bool
	}{
		{name: "discovery endpoint", opts: auth.V3OIDCClientCredentialsOpts{DiscoveryEndpoint: "http://idp.example.com/discovery"}, rejected: true},
		{name: "access token endpoint", opts: auth.V3OIDCClientCredentialsOpts{AccessTokenEndpoint: "http://idp.example.com/token"}, rejected: true},
		{name: "access token endpoint", opts: auth.V3OIDCClientCredentialsOpts{AccessTokenEndpoint: "ftp://idp.example.com/token"}, rejected: true},
		{name: "discovery document token_endpoint", opts: auth.V3OIDCClientCredentialsOpts{DiscoveryEndpoint: "https://example.com/discovery"}, rejected: true},
		// Nothing listens on these loopback ports; the requests are allowed
		// and fail to connect.
		{name: "access token endpoint", opts: auth.V3OIDCClientCredentialsOpts{AccessTokenEndpoint: "http://localhost:1/token"}},
		{name: "access token endpoint", opts: auth.V3OIDCClientCredentialsOpts{AccessTokenEndpoint: "http://[::1]:1/token"}},
		{name: "discovery endpoint", opts: auth.V3OIDCClientCredentialsOpts{DiscoveryEndpoint: "http://127.0.0.1:1/discovery"}},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.opts.AccessTokenEndpoint+tt.opts.DiscoveryEndpoint, func(t *testing.T) {
			opts := tt.opts
			opts.AuthURL = "https://identity.example.com/v3"
			opts.IdentityProviderName, opts.Protocol = "my-idp", "openid"
			opts.ClientID, opts.ClientSecret = "client", "secret"

			_, err := opts.Authenticate(context.Background(), provider)
			th.AssertErr(t, err)
			th.AssertEquals(t, tt.rejected, strings.Contains(err.Error(), tt.name) && strings.Contains(err.Error(), "does not use https"))
		})
	}
	th.AssertEquals(t, int32(0), plaintextRequests.Load())
}

func TestV3OIDCClientCredentialsRedirects(t *testing.T) {
	tests := []struct {
		name      string
		redirect  string
		discovery bool
		err       string
	}{
		{name: "token same origin", redirect: "relative"},
		{name: "token different origin", redirect: "https", err: "token request cannot redirect to another origin"},
		{name: "token downgrade", redirect: "http", err: "request cannot redirect from HTTPS to plaintext"},
		{name: "discovery same origin", redirect: "relative", discovery: true},
		{name: "discovery different origin", redirect: "https", discovery: true},
		{name: "discovery downgrade", redirect: "http", discovery: true, err: "request cannot redirect from HTTPS to plaintext"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var plaintextRequests, otherTokenRequests atomic.Int32
			plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				plaintextRequests.Add(1)
			}))
			defer plain.Close()

			var source *httptest.Server
			other := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					otherTokenRequests.Add(1)
				}
				fmt.Fprintf(w, `{"token_endpoint": %q}`, source.URL+"/token-final")
			}))
			defer other.Close()
			source = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/metadata":
					fmt.Fprintf(w, `{"token_endpoint": %q}`, source.URL+"/token")
				case "/discovery", "/token":
					location := r.URL.Path + "-final"
					switch tt.redirect {
					case "https":
						location = other.URL
					case "http":
						location = plain.URL
					}
					http.Redirect(w, r, location, http.StatusTemporaryRedirect)
				case "/discovery-final":
					fmt.Fprintf(w, `{"token_endpoint": %q}`, source.URL+"/token-final")
				case "/token-final":
					th.TestMethod(t, r, http.MethodPost)
					th.TestHeader(t, r, "Authorization", oidcBasicAuth("my-client-id", "my-client-secret"))
					th.CheckNoErr(t, r.ParseForm())
					th.CheckEquals(t, "client_credentials", r.PostForm.Get("grant_type"))
					fmt.Fprint(w, oidcTokenResponse)
				case oidcFederationPath:
					w.Header().Set("X-Subject-Token", oidcUnscopedToken)
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, oidcKeystoneResponse)
				default:
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer source.Close()

			provider := &gophercloud.ProviderClient{HTTPClient: *source.Client()}
			opts := &auth.V3OIDCClientCredentialsOpts{
				AuthURL:              source.URL,
				IdentityProviderName: "my-idp",
				Protocol:             "openid",
				ClientID:             "my-client-id",
				ClientSecret:         "my-client-secret",
				DiscoveryEndpoint:    source.URL + "/metadata",
			}
			if tt.discovery {
				opts.DiscoveryEndpoint = source.URL + "/discovery"
			}

			_, err := opts.Authenticate(context.Background(), provider)
			if tt.err == "" {
				th.AssertNoErr(t, err)
			} else if err == nil || !strings.Contains(err.Error(), tt.err) {
				t.Fatalf("expected error containing %q, got %v", tt.err, err)
			}
			th.AssertEquals(t, int32(0), plaintextRequests.Load())
			th.AssertEquals(t, int32(0), otherTokenRequests.Load())
			if provider.HTTPClient.CheckRedirect != nil {
				t.Error("provider redirect policy was modified")
			}
		})
	}
}

func TestV3OIDCClientCredentialsRedirectLimits(t *testing.T) {
	var requests atomic.Int32
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		http.Redirect(w, r, "/oauth2/token", http.StatusTemporaryRedirect)
	})

	_, err := oidcOpts(fakeServer).Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
	th.AssertEquals(t, true, strings.Contains(err.Error(), "stopped after 10 redirects"))
	th.AssertEquals(t, int32(10), requests.Load())

	// A redirect policy on the provider's HTTP client replaces the limit.
	want := errors.New("caller rejected redirect")
	provider := &gophercloud.ProviderClient{HTTPClient: http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return want },
	}}
	requests.Store(0)
	_, err = oidcOpts(fakeServer).Authenticate(context.Background(), provider)
	th.AssertEquals(t, true, errors.Is(err, want))
	th.AssertEquals(t, int32(1), requests.Load())
}

func TestV3OIDCClientCredentialsReauthenticates(t *testing.T) {
	var logins atomic.Int32
	var server *httptest.Server
	server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/discovery":
			fmt.Fprintf(w, `{"token_endpoint": %q, "grant_types_supported": ["client_credentials"]}`, server.URL+"/oauth2/token")
		case "/oauth2/token":
			logins.Add(1)
			fmt.Fprint(w, oidcTokenResponse)
		case oidcFederationPath:
			th.CheckEquals(t, "Bearer "+oidcAccessToken, r.Header.Get("Authorization"))
			w.Header().Set("X-Subject-Token", oidcUnscopedToken)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, oidcKeystoneResponse)
		case "/v3/auth/tokens":
			w.Header().Set("X-Subject-Token", fmt.Sprintf("scoped-%d", logins.Load()))
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"token": {"project": {"id": "project"}, "catalog": [{"type": "compute", "endpoints": [{"interface": "public", "url": %q}]}]}}`, server.URL+"/service/")
		case "/service/":
			if r.Header.Get("X-Auth-Token") != "scoped-2" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	opts := &auth.V3OIDCClientCredentialsOpts{
		AuthURL:              server.URL,
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		DiscoveryEndpoint:    server.URL + "/discovery",
		Scope:                &auth.Scope{ProjectID: "project"},
		AllowReauth:          true,
	}
	ctx := context.Background()
	provider, err := config.NewProviderClient(ctx, opts, config.WithHTTPClient(*server.Client()))
	th.AssertNoErr(t, err)
	endpoint, err := provider.EndpointLocator(ctx, gophercloud.EndpointOpts{Type: "compute"})
	th.AssertNoErr(t, err)

	_, err = provider.Request(ctx, http.MethodGet, endpoint, &gophercloud.RequestOpts{OkCodes: []int{http.StatusNoContent}})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, int32(2), logins.Load())
	th.AssertEquals(t, "scoped-2", provider.Token())
}
