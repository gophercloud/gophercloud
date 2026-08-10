package testing

import (
	"context"
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

func startWebSSO(ctx context.Context, client *gophercloud.ServiceClient, opts *websso.AuthOptions) <-chan tokens.CreateResult {
	ready := make(chan struct{})
	opts.BrowserOpener = func(string) error {
		close(ready)
		return nil
	}
	results := make(chan tokens.CreateResult, 1)
	go func() { results <- websso.Authenticate(ctx, client, opts) }()
	<-ready
	return results
}

var _ tokens.AuthOptionsBuilder = (*websso.AuthOptions)(nil)

func TestWebSSOValidationMissingIdentityProvider(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &websso.AuthOptions{
		Protocol: "openid",
	}

	result := websso.Authenticate(context.TODO(), &client, opts)
	if result.Err == nil {
		t.Fatal("Expected error for missing IdentityProviderName")
	}
	if !strings.Contains(result.Err.Error(), "IdentityProviderName") {
		t.Errorf("Expected error about IdentityProviderName, got: %v", result.Err)
	}
}

func TestWebSSOValidationMissingProtocol(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
	}

	result := websso.Authenticate(context.TODO(), &client, opts)
	if result.Err == nil {
		t.Fatal("Expected error for missing Protocol")
	}
	if !strings.Contains(result.Err.Error(), "Protocol") {
		t.Errorf("Expected error about Protocol, got: %v", result.Err)
	}
}

func TestWebSSOWrongOptionsType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}

	opts := &tokens.AuthOptions{}

	result := websso.Authenticate(context.TODO(), &client, opts)
	if result.Err == nil {
		t.Fatal("Expected error for wrong options type")
	}
	if !strings.Contains(result.Err.Error(), "websso.AuthOptions") {
		t.Errorf("Expected error about websso.AuthOptions, got: %v", result.Err)
	}
}

func TestWebSSOScopeMap(t *testing.T) {
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
	}

	scopeMap, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	if scopeMap == nil {
		t.Fatal("Expected non-nil scope map")
	}

	project, ok := scopeMap["project"].(map[string]any)
	if !ok {
		t.Fatalf("Expected project in scope map, got %v", scopeMap)
	}
	projectName, ok := project["name"].(*string)
	if !ok || projectName == nil || *projectName != "my-project" {
		t.Errorf("Expected project name 'my-project', got %v", project["name"])
	}
}

func TestWebSSONoScopeReturnsNil(t *testing.T) {
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
	}

	scopeMap, err := opts.ToTokenV3ScopeMap()
	th.AssertNoErr(t, err)
	if scopeMap != nil {
		t.Errorf("Expected nil scope map for unscoped, got %v", scopeMap)
	}
}

func TestWebSSOCanReauth(t *testing.T) {
	opts := &websso.AuthOptions{AllowReauth: true}
	th.CheckEquals(t, true, opts.CanReauth())

	opts2 := &websso.AuthOptions{AllowReauth: false}
	th.CheckEquals(t, false, opts2.CanReauth())
}

func TestWebSSOCallbackHandlerAcceptsValidPost(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()

	HandleWebSSOTokenValidation(t, fakeKeystone, UnscopedTokenID)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeKeystone.Endpoint(),
	}

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         0,
		Timeout:              10 * time.Second,
	}

	port := findAvailablePort(t)
	opts.RedirectPort = port
	results := startWebSSO(context.Background(), &client, opts)

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	resp, err := http.PostForm(callbackURL, url.Values{"token": {UnscopedTokenID}})
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusOK, resp.StatusCode)

	result := <-results
	th.AssertNoErr(t, result.Err)

	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UnscopedTokenID, token)
}

func TestWebSSOCallbackRejectsWrongMethod(t *testing.T) {
	port := findAvailablePort(t)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "http://localhost/v3/",
	}

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              3 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	startWebSSO(ctx, &client, opts)
	defer cancel()

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	resp, err := http.Get(callbackURL) //nolint:gosec
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestWebSSOCallbackRejectsWrongContentType(t *testing.T) {
	port := findAvailablePort(t)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "http://localhost/v3/",
	}

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              3 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	startWebSSO(ctx, &client, opts)
	defer cancel()

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	resp, err := http.Post(callbackURL, "application/json", strings.NewReader(`{"token":"abc"}`)) //nolint:gosec
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestWebSSOCallbackIgnoresMissingTokenThenAcceptsValidToken(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()

	HandleWebSSOTokenValidation(t, fakeKeystone, UnscopedTokenID)

	port := findAvailablePort(t)

	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeKeystone.Endpoint(),
	}

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              10 * time.Second,
	}

	results := startWebSSO(context.Background(), &client, opts)

	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)

	resp, err := http.PostForm(callbackURL, url.Values{"nottoken": {"abc"}})
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = http.PostForm(callbackURL, url.Values{"token": {UnscopedTokenID}})
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusOK, resp.StatusCode)

	result := <-results
	th.AssertNoErr(t, result.Err)
}

func TestWebSSOCancellationDoesNotWaitForSlowCallback(t *testing.T) {
	port := findAvailablePort(t)
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "http://localhost/v3/",
	}
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	results := startWebSSO(ctx, &client, opts)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	th.AssertNoErr(t, err)

	_, err = fmt.Fprintf(conn, "POST /auth/websso/ HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 1000\r\n\r\ntoken=")
	th.AssertNoErr(t, err)
	cancel()

	select {
	case result := <-results:
		if result.Err == nil || !strings.Contains(result.Err.Error(), "cancelled") {
			t.Fatalf("expected cancellation error, got %v", result.Err)
		}
	case <-time.After(500 * time.Millisecond):
		conn.Close()
		<-results
		t.Fatal("WebSSO cancellation waited for a slow callback connection")
	}
	conn.Close()
}

func TestWebSSOTimeout(t *testing.T) {
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "http://localhost/v3/",
	}

	port := findAvailablePort(t)
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              500 * time.Millisecond,
		BrowserOpener:        func(string) error { return nil },
	}

	result := websso.Authenticate(context.TODO(), &client, opts)
	if result.Err == nil {
		t.Fatal("Expected timeout error")
	}
	if !strings.Contains(result.Err.Error(), "timed out") {
		t.Errorf("Expected timeout error, got: %v", result.Err)
	}
}

func TestWebSSOContextCancellation(t *testing.T) {
	client := gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "http://localhost/v3/",
	}

	port := findAvailablePort(t)
	ctx, cancel := context.WithCancel(context.Background())

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              30 * time.Second,
	}

	results := startWebSSO(ctx, &client, opts)
	cancel()

	res := <-results
	if res.Err == nil {
		t.Fatal("Expected cancellation error")
	}
	if !strings.Contains(res.Err.Error(), "cancel") {
		t.Errorf("Expected cancellation error, got: %v", res.Err)
	}
}
