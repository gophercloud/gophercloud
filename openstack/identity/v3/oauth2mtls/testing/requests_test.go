package testing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oauth2mtls"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func newClient(fakeServer th.FakeServer) *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}
}

func TestAuthenticateV3UsesBearerToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, tokenResponse)
	})
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "X-Subject-Token", tokenID)
		w.Header().Set("X-Subject-Token", tokenID)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, validateTokenResponse)
	})
	fakeServer.Mux.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "Authorization", "Bearer "+tokenID)
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		w.WriteHeader(http.StatusNoContent)
	})

	provider, err := openstack.NewClient(fakeServer.Endpoint() + "v3/")
	th.AssertNoErr(t, err)
	th.AssertNoErr(t, openstack.AuthenticateV3(context.Background(), provider, &oauth2mtls.AuthOptions{
		ClientID: testClientID,
	}, gophercloud.EndpointOpts{}))

	_, err = provider.Request(context.Background(), http.MethodGet, fakeServer.Endpoint()+"service", &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	th.AssertNoErr(t, err)
}

func TestAuthenticateV3Reauthenticates(t *testing.T) {
	const refreshedToken = "refreshed-token"

	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	var tokenRequests atomic.Int32
	fakeServer.Mux.HandleFunc("/v3/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
		token := tokenID
		if tokenRequests.Add(1) > 1 {
			token = refreshedToken
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":%q,"token_type":"Bearer"}`, token)
	})
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Subject-Token")
		th.TestHeader(t, r, "X-Auth-Token", token)
		w.Header().Set("X-Subject-Token", token)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, validateTokenResponse)
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
	th.AssertNoErr(t, openstack.AuthenticateV3(context.Background(), provider, &oauth2mtls.AuthOptions{
		ClientID:    testClientID,
		AllowReauth: true,
	}, gophercloud.EndpointOpts{}))

	_, err = provider.Request(context.Background(), http.MethodGet, fakeServer.Endpoint()+"service", &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, int32(2), tokenRequests.Load())
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleTokenSuccessfully(t, fakeServer)
	handleValidateTokenSuccessfully(t, fakeServer)

	result := oauth2mtls.Create(context.Background(), newClient(fakeServer), &oauth2mtls.AuthOptions{
		ClientID: testClientID,
	})
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, tokenID, token)
	catalog, err := result.ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(catalog.Entries))
}

func TestCreateUsesExplicitOAuth2Endpoint(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/custom/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, tokenResponse)
	})
	handleValidateTokenSuccessfully(t, fakeServer)

	result := oauth2mtls.Create(context.Background(), newClient(fakeServer), &oauth2mtls.AuthOptions{
		OAuth2Endpoint: fakeServer.Endpoint() + "custom/token",
		ClientID:       testClientID,
	})
	th.AssertNoErr(t, result.Err)
}

func TestCreateUsesProviderUserAgent(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	userAgent := "oauth2mtls-test " + gophercloud.DefaultUserAgent
	fakeServer.Mux.HandleFunc("/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "User-Agent", userAgent)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, tokenResponse)
	})
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "User-Agent", userAgent)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, validateTokenResponse)
	})
	client := newClient(fakeServer)
	client.UserAgent.Prepend("oauth2mtls-test")

	result := oauth2mtls.Create(context.Background(), client, &oauth2mtls.AuthOptions{ClientID: testClientID})
	th.AssertNoErr(t, result.Err)
}

func TestCreateRejectsInvalidResponses(t *testing.T) {
	tests := map[string]struct {
		response string
		err      string
	}{
		"invalid JSON":         {response: "not JSON", err: "invalid character"},
		"missing access token": {response: `{"token_type":"Bearer"}`, err: "missing access_token"},
		"unsupported type":     {response: `{"access_token":"token","token_type":"MAC"}`, err: "unsupported token_type"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.response)
			})

			result := oauth2mtls.Create(context.Background(), newClient(fakeServer), &oauth2mtls.AuthOptions{ClientID: testClientID})
			if result.Err == nil || !strings.Contains(result.Err.Error(), tt.err) {
				t.Fatalf("expected error containing %q, got %v", tt.err, result.Err)
			}
		})
	}
}

func TestCreateReturnsResponseErrors(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			})

			result := oauth2mtls.Create(context.Background(), newClient(fakeServer), &oauth2mtls.AuthOptions{ClientID: testClientID})
			if !gophercloud.ResponseCodeIs(result.Err, status) {
				t.Fatalf("expected a %d response error, got %v", status, result.Err)
			}
		})
	}
}

func TestCreateReturnsValidationError(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleTokenSuccessfully(t, fakeServer)
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result := oauth2mtls.Create(context.Background(), newClient(fakeServer), &oauth2mtls.AuthOptions{ClientID: testClientID})
	if !gophercloud.ResponseCodeIs(result.Err, http.StatusUnauthorized) {
		t.Fatalf("expected a 401 response error, got %v", result.Err)
	}
}

func TestCreateValidatesInputs(t *testing.T) {
	client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}}
	result := oauth2mtls.Create(context.Background(), client, &oauth2mtls.AuthOptions{})
	var missing gophercloud.ErrMissingInput
	if !errors.As(result.Err, &missing) {
		t.Fatalf("expected ErrMissingInput, got %v", result.Err)
	}

	var opts *oauth2mtls.AuthOptions
	result = oauth2mtls.Create(context.Background(), client, opts)
	if result.Err == nil {
		t.Fatal("expected typed nil options to fail")
	}

	result = oauth2mtls.Create(context.Background(), nil, &oauth2mtls.AuthOptions{ClientID: testClientID})
	if result.Err == nil {
		t.Fatal("expected nil client to fail")
	}
}

func TestAuthOptions(t *testing.T) {
	opts := &oauth2mtls.AuthOptions{AllowReauth: true}
	th.AssertEquals(t, true, opts.CanReauth())

	scope, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	if scope != nil {
		t.Fatalf("expected nil scope, got %v", scope)
	}
}
