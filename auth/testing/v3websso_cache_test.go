package testing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/auth/tokencache"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

type webSSOMemoryCache map[string]string

func (c webSSOMemoryCache) Get(key string) (string, error) { return c[key], nil }

func (c webSSOMemoryCache) Set(key, value string) error {
	c[key] = value
	return nil
}

func (c webSSOMemoryCache) Delete(key string) error {
	delete(c, key)
	return nil
}

// A cache is configuration, not part of a request.
func (webSSOMemoryCache) MarshalJSON() ([]byte, error) {
	return nil, errors.New("the cache must not be serialized")
}

func newCachedWebSSOOpts(t *testing.T, authURL string, cache tokencache.Cache) *auth.V3WebSSOOpts {
	opts := newWebSSOOpts(t, authURL)
	opts.TokenCache = cache
	opts.CacheNamespace = "profile"
	return opts
}

func TestWebSSOCacheKey(t *testing.T) {
	const endpoint = "https://identity.example.com/v3/"
	opts := auth.V3WebSSOOpts{IdentityProviderName: "my-idp", Protocol: "openid", CacheNamespace: "profile"}
	key := auth.WebSSOCacheKey(endpoint, &opts)

	scoped := opts
	scoped.Scope = &auth.Scope{ProjectID: "project-id"}
	th.AssertEquals(t, key, auth.WebSSOCacheKey(endpoint, &scoped))

	for name, modify := range map[string]func(*auth.V3WebSSOOpts){
		"namespace":         func(o *auth.V3WebSSOOpts) { o.CacheNamespace = "other" },
		"identity provider": func(o *auth.V3WebSSOOpts) { o.IdentityProviderName = "other" },
		"protocol":          func(o *auth.V3WebSSOOpts) { o.Protocol = "saml2" },
	} {
		other := opts
		modify(&other)
		if auth.WebSSOCacheKey(endpoint, &other) == key {
			t.Errorf("the cache key does not depend on the %s", name)
		}
	}
	if auth.WebSSOCacheKey("https://other.example.com/v3/", &opts) == key {
		t.Error("the cache key does not depend on the identity endpoint")
	}
	otherFlow := tokencache.Key(tokencache.KeyOptions{
		Flow: "other", Principal: "profile", IdentityEndpoint: endpoint, IdentityProvider: "my-idp", Protocol: "openid",
	})
	if otherFlow == key {
		t.Error("the cache key does not depend on the flow")
	}
}

func TestV3WebSSOTokenCache(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	handleWebSSOTokens(t, fakeServer, http.StatusOK)
	cache := webSSOMemoryCache{}

	// The first login stores the unscoped token.
	opts := newCachedWebSSOOpts(t, fakeServer.Endpoint(), cache)
	opts.BrowserOpener = postWebSSOToken(t)
	result, err := opts.Authenticate(context.Background(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, webSSOToken, result.TokenID)
	cached, ok := tokencache.Load(cache, auth.WebSSOCacheKey(opts.GetAuthURL(), opts), opts.GetAuthURL())
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, webSSOToken, cached)

	// Later authentications reuse it for any scope, without the browser.
	for _, scope := range []*auth.Scope{nil, {ProjectID: "project-id"}} {
		opts := newCachedWebSSOOpts(t, fakeServer.Endpoint(), cache)
		opts.Scope = scope
		opts.AllowReauth = true
		opts.BrowserOpener = failIfBrowserOpens(t)
		result, err := opts.Authenticate(context.Background(), nil)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, true, result.CanReauth)
		if scope == nil {
			th.AssertEquals(t, webSSOToken, result.TokenID)
		} else {
			th.AssertEquals(t, webSSOScopedToken, result.TokenID)
		}
	}

	// Another namespace has its own login.
	opts = newCachedWebSSOOpts(t, fakeServer.Endpoint(), cache)
	opts.CacheNamespace = "other"
	opened := false
	opts.BrowserOpener = func(target string) error {
		opened = true
		return postWebSSOToken(t)(target)
	}
	_, err = opts.Authenticate(context.Background(), nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, opened)
}

func TestV3WebSSOTokenCacheRevalidatesToken(t *testing.T) {
	browserErr := errors.New("browser unavailable")
	tests := []struct {
		name   string
		status int
		// browser is nil when the browser must not open.
		browser func(string) error
		cached  string
	}{
		{name: "revoked", status: http.StatusUnauthorized, browser: postWebSSOToken(t), cached: webSSOToken},
		{name: "unknown", status: http.StatusNotFound, browser: func(string) error { return browserErr }},
		{name: "Keystone unavailable", status: http.StatusServiceUnavailable, cached: "cached-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)
				w.Header().Set("Content-Type", "application/json")
				if r.Header.Get("X-Subject-Token") == "cached-token" {
					w.WriteHeader(tt.status)
					return
				}
				fmt.Fprint(w, webSSODetails)
			})

			cache := webSSOMemoryCache{}
			opts := newCachedWebSSOOpts(t, fakeServer.Endpoint(), cache)
			key := auth.WebSSOCacheKey(opts.GetAuthURL(), opts)
			tokencache.Persist(cache, key, opts.GetAuthURL(), "cached-token", time.Now().Add(time.Hour))
			opts.BrowserOpener = tt.browser
			if tt.browser == nil {
				opts.BrowserOpener = failIfBrowserOpens(t)
			}

			result, err := opts.Authenticate(context.Background(), nil)
			switch tt.status {
			case http.StatusUnauthorized:
				th.AssertNoErr(t, err)
				th.AssertEquals(t, webSSOToken, result.TokenID)
			case http.StatusNotFound:
				th.AssertEquals(t, true, errors.Is(err, browserErr))
			default:
				th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, tt.status))
			}
			cached, _ := tokencache.Load(cache, key, opts.GetAuthURL())
			th.AssertEquals(t, tt.cached, cached)
		})
	}
}

func TestV3WebSSOTokenCacheRequiresNamespace(t *testing.T) {
	opts := newCachedWebSSOOpts(t, "https://identity.example.com/", webSSOMemoryCache{})
	opts.CacheNamespace = ""
	opts.BrowserOpener = failIfBrowserOpens(t)
	_, err := opts.Authenticate(context.Background(), nil)
	th.AssertErr(t, err)
}
