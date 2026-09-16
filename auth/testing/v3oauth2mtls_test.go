package testing

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const (
	oauth2MTLSToken         = "gAAAAABl-mTLS-token-abc123"
	oauth2MTLSClientID      = "6c3145f4-313d-4910-b3a8-9dfc72da9e75"
	oauth2MTLSTokenResponse = `{"access_token": "gAAAAABl-mTLS-token-abc123", "token_type": "Bearer", "expires_in": 3600}`
	oauth2MTLSTokenDetails  = `{
		"token": {
			"catalog": [{
				"type": "compute",
				"name": "nova",
				"endpoints": [{"interface": "public", "region": "iad1", "url": "http://127.0.0.1:8774/v2.1"}]
			}],
			"expires_at": "2030-06-15T18:00:00.000000Z"
		}
	}`
)

func handleOAuth2MTLSToken(t *testing.T, fakeServer th.FakeServer, path string) {
	fakeServer.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		th.CheckNoErr(t, r.ParseForm())
		th.CheckEquals(t, "client_credentials", r.PostForm.Get("grant_type"))
		th.CheckEquals(t, oauth2MTLSClientID, r.PostForm.Get("client_id"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, oauth2MTLSTokenResponse)
	})
}

func handleOAuth2MTLSTokenDetails(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", oauth2MTLSToken)
		th.TestHeader(t, r, "X-Subject-Token", oauth2MTLSToken)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, oauth2MTLSTokenDetails)
	})
}

func TestV3OAuth2MTLSAuthenticate(t *testing.T) {
	for _, path := range []string{"/v3/OS-OAUTH2/token", "/custom/token"} {
		t.Run(path, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			handleOAuth2MTLSToken(t, fakeServer, path)
			handleOAuth2MTLSTokenDetails(t, fakeServer)

			opts := &auth.V3OAuth2MTLSOpts{
				AuthURL:     fakeServer.Endpoint(),
				ClientID:    oauth2MTLSClientID,
				AllowReauth: true,
			}
			if path == "/custom/token" {
				opts.OAuth2Endpoint = fakeServer.Endpoint() + "custom/token"
			}

			result, err := opts.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, oauth2MTLSToken, result.TokenID)
			th.AssertEquals(t, "Bearer", result.TokenType)
			th.AssertEquals(t, true, result.CanReauth)
			th.AssertEquals(t, 1, len(result.Catalog.Entries))
		})
	}
}

func TestV3OAuth2MTLSRejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name          string
		tokenStatus   int
		tokenResponse string
		detailsStatus int
		err           string
	}{
		{name: "invalid JSON", tokenStatus: http.StatusOK, tokenResponse: "not JSON", err: "invalid character"},
		{name: "missing access token", tokenStatus: http.StatusOK, tokenResponse: `{"token_type":"Bearer"}`, err: "missing access_token"},
		{name: "unsupported token type", tokenStatus: http.StatusOK, tokenResponse: `{"access_token":"token","token_type":"MAC"}`, err: "unsupported token_type"},
		{name: "bad request", tokenStatus: http.StatusBadRequest},
		{name: "unauthorized", tokenStatus: http.StatusUnauthorized},
		{name: "rejected token", tokenStatus: http.StatusOK, tokenResponse: oauth2MTLSTokenResponse, detailsStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/v3/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.tokenStatus)
				fmt.Fprint(w, tt.tokenResponse)
			})
			fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.detailsStatus)
			})

			opts := &auth.V3OAuth2MTLSOpts{AuthURL: fakeServer.Endpoint(), ClientID: oauth2MTLSClientID}
			_, err := opts.Authenticate(context.Background(), nil)
			switch {
			case tt.err != "":
				if err == nil || !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("expected error containing %q, got %v", tt.err, err)
				}
			case tt.detailsStatus != 0:
				th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, tt.detailsStatus))
			default:
				th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, tt.tokenStatus))
			}
		})
	}
}

func TestV3OAuth2MTLSRequiresInput(t *testing.T) {
	var missing gophercloud.ErrMissingInput

	_, err := (&auth.V3OAuth2MTLSOpts{ClientID: oauth2MTLSClientID}).Authenticate(context.Background(), nil)
	th.AssertEquals(t, true, errors.As(err, &missing))
	th.AssertEquals(t, "AuthURL", missing.Argument)

	_, err = (&auth.V3OAuth2MTLSOpts{AuthURL: "https://identity.example.com/v3"}).Authenticate(context.Background(), nil)
	th.AssertEquals(t, true, errors.As(err, &missing))
	th.AssertEquals(t, "ClientID", missing.Argument)

	var opts *auth.V3OAuth2MTLSOpts
	_, err = opts.Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
}

func TestV3OAuth2MTLSProviderUsesCertificateAndBearerToken(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if len(r.TLS.PeerCertificates) == 0 {
			t.Error("missing client certificate")
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v3/OS-OAUTH2/token":
			fmt.Fprint(w, oauth2MTLSTokenResponse)
		case "/v3/auth/tokens":
			th.CheckEquals(t, oauth2MTLSToken, r.Header.Get("X-Subject-Token"))
			fmt.Fprint(w, oauth2MTLSTokenDetails)
		case "/service":
			th.CheckEquals(t, "Bearer "+oauth2MTLSToken, r.Header.Get("Authorization"))
			th.CheckEquals(t, oauth2MTLSToken, r.Header.Get("X-Auth-Token"))
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	server.TLS = &tls.Config{ClientAuth: tls.RequireAnyClientCert}
	server.StartTLS()
	defer server.Close()

	transport := server.Client().Transport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	transport.TLSClientConfig.Certificates = server.TLS.Certificates

	provider, err := config.NewProviderClient(context.Background(), &auth.V3OAuth2MTLSOpts{
		AuthURL:  server.URL,
		ClientID: oauth2MTLSClientID,
	}, config.WithHTTPClient(http.Client{Transport: transport}))
	th.AssertNoErr(t, err)

	_, err = provider.Request(context.Background(), http.MethodGet, server.URL+"/service", &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, int32(3), requests.Load())
}

func TestV3OAuth2MTLSReauthenticates(t *testing.T) {
	const refreshedToken = "refreshed-token"

	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	var tokenRequests atomic.Int32
	fakeServer.Mux.HandleFunc("/v3/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
		token := oauth2MTLSToken
		if tokenRequests.Add(1) > 1 {
			token = refreshedToken
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":%q,"token_type":"Bearer"}`, token)
	})
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", r.Header.Get("X-Subject-Token"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, oauth2MTLSTokenDetails)
	})
	fakeServer.Mux.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+refreshedToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	provider, err := openstack.NewClient(fakeServer.Endpoint() + "v3/")
	th.AssertNoErr(t, err)
	th.AssertNoErr(t, openstack.Authenticate(context.Background(), provider, &auth.V3OAuth2MTLSOpts{
		AuthURL:     fakeServer.Endpoint(),
		ClientID:    oauth2MTLSClientID,
		AllowReauth: true,
	}))

	_, err = provider.Request(context.Background(), http.MethodGet, fakeServer.Endpoint()+"service", &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, int32(2), tokenRequests.Load())
}
