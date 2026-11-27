package testing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const (
	webSSOToken       = "unscoped-websso-token"
	webSSOScopedToken = "scoped-websso-token"
	webSSODetails     = `{"token": {"methods": ["mapped"], "expires_at": "2099-01-01T00:00:00Z", "user": {"id": "user-id"}}}`
	webSSOScoped      = `{
		"token": {
			"methods": ["token"],
			"expires_at": "2099-01-01T00:00:00Z",
			"project": {"id": "project-id"},
			"catalog": [{"type": "compute", "endpoints": [{"interface": "public", "url": "https://compute.example/"}]}]
		}
	}`
)

// handleWebSSOTokens serves the details of the unscoped token with
// detailsStatus, and exchanges it for a token scoped to project-id.
func handleWebSSOTokens(t *testing.T, fakeServer th.FakeServer, detailsStatus int) {
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			th.TestHeader(t, r, "X-Auth-Token", webSSOToken)
			th.TestHeader(t, r, "X-Subject-Token", webSSOToken)
			w.WriteHeader(detailsStatus)
			fmt.Fprint(w, webSSODetails)
			return
		}
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		th.TestJSONRequest(t, r, fmt.Sprintf(`{"auth": {"identity": {"methods": ["token"], "token": {"id": %q}}, "scope": {"project": {"id": "project-id"}}}}`, webSSOToken))
		w.Header().Set("X-Subject-Token", webSSOScopedToken)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, webSSOScoped)
	})
}

func newWebSSOOpts(t *testing.T, authURL string) *auth.V3WebSSOOpts {
	return &auth.V3WebSSOOpts{
		AuthURL:              authURL,
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "127.0.0.1",
		RedirectPort:         freeLoopbackPort(t),
		Timeout:              10 * time.Second,
	}
}

func freeLoopbackPort(t *testing.T) int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	th.AssertNoErr(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// webSSOCallback returns the request a browser sends when Keystone's WebSSO
// page posts the token to target.
func webSSOCallback(t *testing.T, target string, form url.Values) *http.Request {
	request, err := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	th.AssertNoErr(t, err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Sec-Fetch-Site", "cross-site")
	request.Header.Set("Origin", "null")
	return request
}

func sendWebSSOCallback(t *testing.T, request *http.Request) (int, string) {
	response, err := http.DefaultClient.Do(request)
	th.AssertNoErr(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	th.AssertNoErr(t, err)
	return response.StatusCode, string(body)
}

// postWebSSOToken returns a BrowserOpener that completes the login by posting
// the token to the callback named in the WebSSO URL.
func postWebSSOToken(t *testing.T) func(string) error {
	return func(target string) error {
		parsed, err := url.Parse(target)
		th.AssertNoErr(t, err)
		status, _ := sendWebSSOCallback(t, webSSOCallback(t, parsed.Query().Get("origin"), url.Values{"token": {webSSOToken}}))
		th.AssertEquals(t, http.StatusOK, status)
		return nil
	}
}

func failIfBrowserOpens(t *testing.T) func(string) error {
	return func(string) error {
		t.Error("unexpected browser")
		return errors.New("unexpected browser")
	}
}

func TestV3WebSSOAuthenticate(t *testing.T) {
	tests := []struct {
		name          string
		detailsStatus int
		scope         *auth.Scope
		token         string
	}{
		{name: "unscoped", detailsStatus: http.StatusOK, token: webSSOToken},
		{name: "unscoped through a proxy", detailsStatus: http.StatusNonAuthoritativeInfo, token: webSSOToken},
		{name: "scoped", detailsStatus: http.StatusOK, scope: &auth.Scope{ProjectID: "project-id"}, token: webSSOScopedToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			handleWebSSOTokens(t, fakeServer, tt.detailsStatus)

			opts := newWebSSOOpts(t, fakeServer.Endpoint())
			opts.Scope = tt.scope
			opts.AllowReauth = true
			opts.BrowserOpener = func(target string) error {
				parsed, err := url.Parse(target)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, "/v3/auth/OS-FEDERATION/identity_providers/my-idp/protocols/openid/websso", parsed.Path)
				origin := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", opts.RedirectPort)
				th.AssertEquals(t, origin, parsed.Query().Get("origin"))

				status, page := sendWebSSOCallback(t, webSSOCallback(t, origin, url.Values{"token": {webSSOToken}}))
				th.AssertEquals(t, http.StatusOK, status)
				th.AssertEquals(t, true, strings.Contains(page, "window.close()"))

				// Only the first token is accepted.
				status, _ = sendWebSSOCallback(t, webSSOCallback(t, origin, url.Values{"token": {"another-token"}}))
				th.AssertEquals(t, http.StatusConflict, status)
				return nil
			}

			provider := &gophercloud.ProviderClient{TokenID: "provider-token"}
			result, err := opts.Authenticate(context.Background(), provider)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, tt.token, result.TokenID)
			th.AssertEquals(t, true, result.CanReauth)
			if tt.scope != nil {
				th.AssertEquals(t, "project-id", result.Project.ID)
				endpoint, err := result.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "compute"})
				th.AssertNoErr(t, err)
				th.AssertEquals(t, "https://compute.example/", endpoint)
			}
		})
	}
}

func TestV3WebSSOProvider(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleWebSSOTokens(t, fakeServer, http.StatusOK)

	opts := newWebSSOOpts(t, fakeServer.Endpoint())
	opts.RedirectHost = ""
	opts.Scope = &auth.Scope{ProjectID: "project-id"}
	opts.BrowserOpener = func(target string) error {
		th.AssertEquals(t, true, strings.Contains(target, url.QueryEscape(fmt.Sprintf("http://localhost:%d/", opts.RedirectPort))))
		return postWebSSOToken(t)(target)
	}

	provider, err := openstack.NewClient(fakeServer.Endpoint())
	th.AssertNoErr(t, err)
	th.AssertNoErr(t, openstack.Authenticate(context.Background(), provider, opts))
	th.AssertEquals(t, webSSOScopedToken, provider.Token())
}

func TestV3WebSSORequiresValidInput(t *testing.T) {
	var missing gophercloud.ErrMissingInput
	var scope gophercloud.ErrScopeProjectDomainIDOrProjectDomainName
	tests := []struct {
		name   string
		modify func(*auth.V3WebSSOOpts)
		target any
		err    string
	}{
		{name: "auth URL", modify: func(o *auth.V3WebSSOOpts) { o.AuthURL = "" }, target: &missing, err: "AuthURL"},
		{name: "identity provider", modify: func(o *auth.V3WebSSOOpts) { o.IdentityProviderName = "" }, target: &missing, err: "IdentityProviderName"},
		{name: "protocol", modify: func(o *auth.V3WebSSOOpts) { o.Protocol = "" }, target: &missing, err: "Protocol"},
		{name: "negative port", modify: func(o *auth.V3WebSSOOpts) { o.RedirectPort = -1 }, err: "RedirectPort"},
		{name: "port too large", modify: func(o *auth.V3WebSSOOpts) { o.RedirectPort = 65536 }, err: "RedirectPort"},
		{name: "negative timeout", modify: func(o *auth.V3WebSSOOpts) { o.Timeout = -time.Second }, err: "Timeout"},
		{name: "non-loopback host", modify: func(o *auth.V3WebSSOOpts) { o.RedirectHost = "0.0.0.0" }, err: "RedirectHost"},
		{name: "host name", modify: func(o *auth.V3WebSSOOpts) { o.RedirectHost = "example.com" }, err: "RedirectHost"},
		{name: "scope", modify: func(o *auth.V3WebSSOOpts) { o.Scope = &auth.Scope{ProjectName: "project"} }, target: &scope},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := newWebSSOOpts(t, "https://identity.example.com/")
			opts.BrowserOpener = failIfBrowserOpens(t)
			tt.modify(opts)

			_, err := opts.Authenticate(context.Background(), nil)
			th.AssertErr(t, err)
			if tt.target != nil {
				th.AssertEquals(t, true, errors.As(err, tt.target))
			}
			th.AssertEquals(t, true, strings.Contains(err.Error(), tt.err))
		})
	}

	var opts *auth.V3WebSSOOpts
	_, err := opts.Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
}

func TestV3WebSSOCallbackRejectsInvalidRequests(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleWebSSOTokens(t, fakeServer, http.StatusOK)
	keystone := strings.TrimSuffix(fakeServer.Endpoint(), "/")

	token := url.Values{"token": {webSSOToken}}
	tests := []struct {
		name   string
		status int
		modify func(*http.Request)
	}{
		{name: "method", status: http.StatusMethodNotAllowed, modify: func(r *http.Request) { r.Method = http.MethodGet }},
		{name: "path", status: http.StatusNotFound, modify: func(r *http.Request) { r.URL.Path += "unexpected" }},
		{name: "content type", status: http.StatusUnsupportedMediaType, modify: func(r *http.Request) { r.Header.Set("Content-Type", "application/json") }},
		{name: "malformed content type", status: http.StatusUnsupportedMediaType, modify: func(r *http.Request) { r.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset") }},
		{name: "fetch mode", status: http.StatusBadRequest, modify: func(r *http.Request) { r.Header.Del("Sec-Fetch-Mode") }},
		{name: "fetch destination", status: http.StatusBadRequest, modify: func(r *http.Request) { r.Header.Set("Sec-Fetch-Dest", "empty") }},
		{name: "direct navigation", status: http.StatusBadRequest, modify: func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "none") }},
		{name: "origin", status: http.StatusBadRequest, modify: func(r *http.Request) { r.Header.Set("Origin", "https://attacker.example") }},
		{name: "referer", status: http.StatusBadRequest, modify: func(r *http.Request) { r.Header.Set("Referer", "https://attacker.example/") }},
		{name: "token in the query", status: http.StatusBadRequest, modify: func(r *http.Request) {
			r.URL.RawQuery = token.Encode()
			r.Body = http.NoBody
		}},
		{name: "oversized body", status: http.StatusBadRequest, modify: func(r *http.Request) {
			body := "padding=" + strings.Repeat("x", 64<<10) + "&" + token.Encode()
			r.Body = io.NopCloser(strings.NewReader(body))
			r.ContentLength = int64(len(body))
		}},
		{name: "Keystone origin and referer", status: http.StatusOK, modify: func(r *http.Request) {
			r.Header.Set("Origin", keystone)
			r.Header.Set("Referer", keystone+"/v3/auth/OS-FEDERATION/websso/openid")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := newWebSSOOpts(t, fakeServer.Endpoint())
			opts.BrowserOpener = func(target string) error {
				parsed, err := url.Parse(target)
				th.AssertNoErr(t, err)
				request := webSSOCallback(t, parsed.Query().Get("origin"), token)
				tt.modify(request)
				status, _ := sendWebSSOCallback(t, request)
				th.AssertEquals(t, tt.status, status)
				if tt.status == http.StatusOK {
					return nil
				}
				// A rejected request does not end the login.
				return postWebSSOToken(t)(target)
			}

			result, err := opts.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, webSSOToken, result.TokenID)
		})
	}
}

func TestV3WebSSOListensOnBothLoopbackAddresses(t *testing.T) {
	probe, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Skip("IPv6 loopback is unavailable")
	}
	probe.Close()

	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleWebSSOTokens(t, fakeServer, http.StatusOK)

	for _, host := range []string{"127.0.0.1", "[::1]"} {
		opts := newWebSSOOpts(t, fakeServer.Endpoint())
		opts.RedirectHost = "localhost"
		opts.BrowserOpener = func(string) error {
			callback := fmt.Sprintf("http://%s:%d/auth/websso/", host, opts.RedirectPort)
			status, _ := sendWebSSOCallback(t, webSSOCallback(t, callback, url.Values{"token": {webSSOToken}}))
			th.AssertEquals(t, http.StatusOK, status)
			return nil
		}
		_, err := opts.Authenticate(context.Background(), nil)
		th.AssertNoErr(t, err)
	}

	// Another process on one of the addresses could receive the token.
	opts := newWebSSOOpts(t, fakeServer.Endpoint())
	opts.RedirectHost = "localhost"
	opts.BrowserOpener = failIfBrowserOpens(t)
	other, err := net.Listen("tcp", net.JoinHostPort("::1", strconv.Itoa(opts.RedirectPort)))
	th.AssertNoErr(t, err)
	defer other.Close()
	_, err = opts.Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
	th.AssertEquals(t, true, strings.Contains(err.Error(), "callback server"))
}

func TestV3WebSSOReturnsErrors(t *testing.T) {
	t.Run("browser", func(t *testing.T) {
		browserErr := errors.New("browser unavailable")
		opts := &auth.V3WebSSOOpts{AuthURL: "https://identity.example.com/", IdentityProviderName: "my-idp", Protocol: "openid"}
		opts.BrowserOpener = func(target string) error {
			parsed, err := url.Parse(target)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, "http://localhost:9990/auth/websso/", parsed.Query().Get("origin"))
			return browserErr
		}
		_, err := opts.Authenticate(context.Background(), nil)
		if err != nil && strings.Contains(err.Error(), "callback server") {
			t.Skip("the default callback port is in use")
		}
		th.AssertEquals(t, true, errors.Is(err, browserErr))
	})

	t.Run("timeout", func(t *testing.T) {
		opts := newWebSSOOpts(t, "https://identity.example.com/")
		opts.Timeout = 50 * time.Millisecond
		opts.BrowserOpener = func(string) error { return nil }
		_, err := opts.Authenticate(context.Background(), nil)
		th.AssertErr(t, err)
		th.AssertEquals(t, true, strings.Contains(err.Error(), "timed out"))
	})

	t.Run("cancelled before the browser opens", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		opts := newWebSSOOpts(t, "https://identity.example.com/")
		opts.BrowserOpener = failIfBrowserOpens(t)
		_, err := opts.Authenticate(ctx, nil)
		th.AssertEquals(t, true, errors.Is(err, context.Canceled))
	})

	t.Run("cancelled during a slow callback", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		opts := newWebSSOOpts(t, "https://identity.example.com/")
		opts.BrowserOpener = func(string) error {
			conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(opts.RedirectPort)))
			th.AssertNoErr(t, err)
			t.Cleanup(func() { conn.Close() })
			request := webSSOCallback(t, "http://localhost/auth/websso/", nil)
			request.ContentLength = 1000
			request.Body = io.NopCloser(strings.NewReader("token="))
			go func() { _ = request.Write(conn) }()
			time.AfterFunc(50*time.Millisecond, cancel)
			return nil
		}

		start := time.Now()
		_, err := opts.Authenticate(ctx, nil)
		th.AssertEquals(t, true, errors.Is(err, context.Canceled))
		if elapsed := time.Since(start); elapsed > 2*time.Second {
			t.Fatalf("cancellation waited %s for the callback connection", elapsed)
		}
	})

	t.Run("rejected token", func(t *testing.T) {
		fakeServer := th.SetupHTTP()
		defer fakeServer.Teardown()
		handleWebSSOTokens(t, fakeServer, http.StatusInternalServerError)

		opts := newWebSSOOpts(t, fakeServer.Endpoint())
		opts.BrowserOpener = postWebSSOToken(t)
		provider := &gophercloud.ProviderClient{TokenID: "provider-token"}
		_, err := opts.Authenticate(context.Background(), provider)
		th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusInternalServerError))
		th.AssertEquals(t, "provider-token", provider.TokenID)
	})
}
