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

func TestWebSSOListenerStartsBeforeBrowser(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()
	handleWebSSOTokenValidation(t, fakeKeystone, testToken)

	result := websso.Authenticate(context.Background(), serviceClient(fakeKeystone.Endpoint()), &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "127.0.0.1",
		RedirectPort:         findAvailablePort(t),
		BrowserOpener: func(target string) error {
			parsed, err := url.Parse(target)
			if err != nil {
				return err
			}
			request := newWebSSOCallbackRequest(t, parsed.Query().Get("origin"), url.Values{"token": {testToken}})
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected callback status: %s", response.Status)
			}
			return nil
		},
	})

	th.AssertNoErr(t, result.Err)
	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, testToken, token)
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

func TestWebSSOCallbackRejectsInvalidRequest(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleWebSSOTokenValidation(t, fakeServer, testToken)
	tests := []struct {
		name       string
		path       string
		statusCode int
		modify     func(*http.Request)
	}{
		{
			name:       "malformed content type",
			statusCode: http.StatusUnsupportedMediaType,
			modify: func(request *http.Request) {
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded-malformed")
			},
		},
		{
			name:       "unexpected path",
			path:       "unexpected",
			statusCode: http.StatusNotFound,
		},
		{
			name: "missing fetch mode",
			modify: func(request *http.Request) {
				request.Header.Del("Sec-Fetch-Mode")
			},
		},
		{
			name: "unexpected fetch destination",
			modify: func(request *http.Request) {
				request.Header.Set("Sec-Fetch-Dest", "empty")
			},
		},
		{
			name: "direct navigation",
			modify: func(request *http.Request) {
				request.Header.Set("Sec-Fetch-Site", "none")
			},
		},
		{
			name: "unexpected origin",
			modify: func(request *http.Request) {
				request.Header.Set("Origin", "https://attacker.example")
			},
		},
		{
			name: "unexpected referer",
			modify: func(request *http.Request) {
				request.Header.Set("Referer", "https://attacker.example/")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := findAvailablePort(t)
			results := startWebSSO(t, context.Background(), serviceClient(fakeServer.Endpoint()), &websso.AuthOptions{
				IdentityProviderName: "my-idp",
				Protocol:             "openid",
				RedirectHost:         "127.0.0.1",
				RedirectPort:         port,
				Timeout:              5 * time.Second,
			})

			callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
			request := newWebSSOCallbackRequest(t, callbackURL+test.path, url.Values{"token": {testToken}})
			if test.modify != nil {
				test.modify(request)
			}
			response, err := http.DefaultClient.Do(request)
			th.AssertNoErr(t, err)
			response.Body.Close()
			wantStatus := test.statusCode
			if wantStatus == 0 {
				wantStatus = http.StatusBadRequest
			}
			if response.StatusCode != wantStatus {
				t.Fatalf("unexpected callback status: got %d, want %d", response.StatusCode, wantStatus)
			}

			response, err = http.DefaultClient.Do(newWebSSOCallbackRequest(t, callbackURL, url.Values{"token": {testToken}}))
			th.AssertNoErr(t, err)
			response.Body.Close()
			if response.StatusCode != http.StatusOK {
				t.Fatalf("unexpected valid callback status: %d", response.StatusCode)
			}

			result := <-results
			th.AssertNoErr(t, result.Err)
		})
	}
}

func newWebSSOCallbackRequest(t *testing.T, target string, values url.Values) *http.Request {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode()))
	th.AssertNoErr(t, err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Sec-Fetch-Site", "cross-site")
	request.Header.Set("Origin", "null")
	return request
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
