package testing

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const testToken = "unscoped-websso-token"

var _ tokens.AuthOptionsBuilder = (*websso.AuthOptions)(nil)

func TestAuthenticateRejectsTypedNilOptions(t *testing.T) {
	client := serviceClient("http://identity.example/v3/")
	var opts *websso.AuthOptions
	result := websso.Authenticate(context.Background(), client, opts)
	if result.Err == nil || !strings.Contains(result.Err.Error(), "non-nil") {
		t.Fatalf("expected typed nil options error, got %v", result.Err)
	}
}

func TestBrowserOpenerFailureIsReturned(t *testing.T) {
	port := findAvailablePort(t)
	client := serviceClient("http://identity.example/v3/")
	wantErr := errors.New("browser unavailable")
	var openedURL string

	result := websso.Authenticate(context.Background(), client, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "127.0.0.1",
		RedirectPort:         port,
		BrowserOpener: func(target string) error {
			openedURL = target
			return wantErr
		},
	})

	if !errors.Is(result.Err, wantErr) {
		t.Fatalf("expected browser opener error, got %v", result.Err)
	}
	parsed, err := url.Parse(openedURL)
	th.AssertNoErr(t, err)
	wantOrigin := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	if got := parsed.Query().Get("origin"); got != wantOrigin {
		t.Fatalf("unexpected WebSSO origin: got %q, want %q", got, wantOrigin)
	}
}

func TestCacheRequiresExplicitNamespace(t *testing.T) {
	client := serviceClient("http://identity.example/v3/")
	result := websso.Authenticate(context.Background(), client, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		TokenCache:           newMemoryCache(),
	})
	if result.Err == nil || !strings.Contains(result.Err.Error(), "CacheNamespace") {
		t.Fatalf("expected CacheNamespace validation error, got %v", result.Err)
	}
}

func TestRedirectHostMustBeLoopback(t *testing.T) {
	client := serviceClient("http://identity.example/v3/")
	result := websso.Authenticate(context.Background(), client, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "0.0.0.0",
	})
	if result.Err == nil || !strings.Contains(result.Err.Error(), "loopback") {
		t.Fatalf("expected loopback validation error, got %v", result.Err)
	}
}

func TestInvalidScopeDoesNotOpenBrowser(t *testing.T) {
	client := serviceClient("http://identity.example/v3/")
	opened := false
	result := websso.Authenticate(context.Background(), client, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		Scope: tokens.Scope{
			ProjectName: "project",
		},
		RedirectPort: findAvailablePort(t),
		BrowserOpener: func(string) error {
			opened = true
			return errors.New("browser should not open")
		},
	})

	if result.Err == nil || !strings.Contains(result.Err.Error(), "DomainID or DomainName") {
		t.Fatalf("expected invalid scope error, got %v", result.Err)
	}
	if opened {
		t.Fatal("browser opened before scope validation")
	}
}

func TestMalformedContentTypeDoesNotConsumeCallback(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Subject-Token", testToken)
		w.Header().Set("X-Subject-Token", testToken)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"token":{"expires_at":"2035-01-01T00:00:00Z","methods":["mapped"],"user":{"id":"user-id","name":"user","domain":{"id":"default","name":"Default"}}}}`)
	})

	port := findAvailablePort(t)
	client := serviceClient(fakeServer.Endpoint())
	results := startWebSSO(context.Background(), client, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "127.0.0.1",
		RedirectPort:         port,
		Timeout:              5 * time.Second,
	})

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	request, err := http.NewRequest(http.MethodPost, callbackURL, strings.NewReader("token="+testToken))
	th.AssertNoErr(t, err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded-malformed")
	response, err := http.DefaultClient.Do(request)
	th.AssertNoErr(t, err)
	response.Body.Close()
	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("unexpected malformed callback status: %d", response.StatusCode)
	}

	response, err = http.PostForm(callbackURL, url.Values{"token": {testToken}})
	th.AssertNoErr(t, err)
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected valid callback status: %d", response.StatusCode)
	}

	result := <-results
	th.AssertNoErr(t, result.Err)
	actualToken, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	if actualToken != testToken {
		t.Fatalf("unexpected token: got %q, want %q", actualToken, testToken)
	}
}

func serviceClient(endpoint string) *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       endpoint,
	}
}

func findAvailablePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	th.AssertNoErr(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
