package testing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokencache"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

type memoryCache map[string]string

func newMemoryCache() memoryCache { return make(memoryCache) }

func (c memoryCache) Get(key string) (string, error) { return c[key], nil }
func (c memoryCache) Set(key, value string) error {
	c[key] = value
	return nil
}
func (c memoryCache) Delete(key string) error {
	delete(c, key)
	return nil
}

func TestWebSSOCacheKeySeparatesFlows(t *testing.T) {
	endpoint := "https://keystone.example.com/v3"
	webSSOKey := websso.CacheKey(endpoint, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		CacheNamespace:       "profile",
	})
	otherKey := tokencache.Key(tokencache.KeyOptions{
		Flow:             "other-flow",
		IdentityEndpoint: endpoint,
		Principal:        "profile",
	})
	if webSSOKey == otherKey {
		t.Fatalf("authentication flows must not share cache key %q", webSSOKey)
	}
}

func TestWebSSOCacheKeySeparatesNamespaces(t *testing.T) {
	endpoint := "https://keystone.example.com/v3"
	base := websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		CacheNamespace:       "account-a",
	}
	other := base
	other.CacheNamespace = "account-b"

	if first, second := websso.CacheKey(endpoint, &base), websso.CacheKey(endpoint, &other); first == second {
		t.Fatalf("different cache namespaces share key %q", first)
	}
}

func TestWebSSOCacheKeyIgnoresScope(t *testing.T) {
	endpoint := "https://keystone.example.com/v3"
	base := websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		CacheNamespace:       "profile",
	}
	projectA := base
	projectA.Scope = tokens.Scope{ProjectID: "project-a"}
	projectB := base
	projectB.Scope = tokens.Scope{ProjectID: "project-b"}

	if first, second := websso.CacheKey(endpoint, &projectA), websso.CacheKey(endpoint, &projectB); first != second {
		t.Fatalf("scope changed unscoped token cache key: %q != %q", first, second)
	}
}

func TestWebSSOCacheReusesUnscopedTokenAcrossProjects(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()

	const (
		unscopedToken = "unscoped-token"
		projectAToken = "project-a-token"
		projectBToken = "project-b-token"
	)
	var scopeRequests atomic.Int32
	fakeKeystone.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			th.TestHeader(t, r, "X-Subject-Token", unscopedToken)
			w.Header().Set("X-Subject-Token", unscopedToken)
			fmt.Fprint(w, tokenResponse(""))
		case http.MethodPost:
			var project, token string
			switch scopeRequests.Add(1) {
			case 1:
				project, token = "project-a", projectAToken
			case 2:
				project, token = "project-b", projectBToken
			default:
				t.Fatal("unexpected extra scope request")
			}
			th.TestJSONRequest(t, r, fmt.Sprintf(`{"auth":{"identity":{"methods":["token"],"token":{"id":%q}},"scope":{"project":{"id":%q}}}}`, unscopedToken, project))
			w.Header().Set("X-Subject-Token", token)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, tokenResponse(project))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	cache := newMemoryCache()
	client := gophercloud.ServiceClient{
		ProviderClient: newTestProviderClient(),
		Endpoint:       fakeKeystone.Endpoint() + "v3/",
	}
	base := websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		Timeout:              10 * time.Second,
		TokenCache:           cache,
		CacheNamespace:       "profile",
	}

	projectA := base
	projectA.Scope = tokens.Scope{ProjectID: "project-a"}
	projectA.RedirectPort = findAvailablePort(t)
	results := startWebSSO(context.Background(), &client, &projectA)
	response, err := http.PostForm(fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", projectA.RedirectPort), url.Values{"token": {unscopedToken}})
	th.AssertNoErr(t, err)
	response.Body.Close()
	result := <-results
	th.AssertNoErr(t, result.Err)
	token, err := result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, projectAToken, token)

	projectB := base
	projectB.Scope = tokens.Scope{ProjectID: "project-b"}
	projectB.RedirectPort = findAvailablePort(t)
	projectB.BrowserOpener = func(string) error { return errors.New("browser opened on cache hit") }
	result = websso.Authenticate(context.Background(), &client, &projectB)
	th.AssertNoErr(t, result.Err)
	token, err = result.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, projectBToken, token)
	th.CheckEquals(t, int32(2), scopeRequests.Load())
}

func TestWebSSOKeepsKnownTokenIDAfterValidation(t *testing.T) {
	const (
		unscopedToken = "unscoped-token"
		scopedToken   = "scoped-token"
	)

	tests := []struct {
		name   string
		cached bool
		scoped bool
	}{
		{name: "browser unscoped"},
		{name: "browser scoped", scoped: true},
		{name: "cache scoped", cached: true, scoped: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeKeystone := th.SetupHTTP()
			defer fakeKeystone.Teardown()

			fakeKeystone.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					th.TestHeader(t, r, "X-Subject-Token", unscopedToken)
					fmt.Fprint(w, tokenResponse(""))
				case http.MethodPost:
					if !tt.scoped {
						t.Fatal("unexpected scope request")
					}
					th.TestJSONRequest(t, r, fmt.Sprintf(`{"auth":{"identity":{"methods":["token"],"token":{"id":%q}},"scope":{"project":{"id":"project-id"}}}}`, unscopedToken))
					w.Header().Set("X-Subject-Token", scopedToken)
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, tokenResponse("project-id"))
				default:
					t.Fatalf("unexpected method %s", r.Method)
				}
			})

			cache := newMemoryCache()
			client := gophercloud.ServiceClient{
				ProviderClient: newTestProviderClient(),
				Endpoint:       fakeKeystone.Endpoint() + "v3/",
			}
			opts := &websso.AuthOptions{
				IdentityProviderName: "my-idp",
				Protocol:             "openid",
				RedirectPort:         findAvailablePort(t),
				Timeout:              10 * time.Second,
			}
			if tt.scoped {
				opts.Scope = tokens.Scope{ProjectID: "project-id"}
			}
			if tt.cached {
				opts.TokenCache = cache
				opts.CacheNamespace = "profile"
				key := websso.CacheKey(client.Endpoint, opts)
				cached := fmt.Sprintf(`{"token_id":%q,"expires_at":%q,"endpoint":%q}`, unscopedToken, time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano), client.Endpoint)
				th.AssertNoErr(t, cache.Set(key, cached))
				opts.BrowserOpener = func(string) error { return errors.New("browser opened on cache hit") }
			}

			var result tokens.CreateResult
			if tt.cached {
				result = websso.Authenticate(context.Background(), &client, opts)
			} else {
				results := startWebSSO(context.Background(), &client, opts)
				response, err := http.PostForm(fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", opts.RedirectPort), url.Values{"token": {unscopedToken}})
				th.AssertNoErr(t, err)
				response.Body.Close()
				result = <-results
			}

			th.AssertNoErr(t, result.Err)
			tokenID, err := result.ExtractTokenID()
			th.AssertNoErr(t, err)
			wantToken := unscopedToken
			if tt.scoped {
				wantToken = scopedToken
			}
			th.CheckEquals(t, wantToken, tokenID)
		})
	}
}

func TestWebSSOCacheValidationCancellationDoesNotOpenBrowser(t *testing.T) {
	cache := newMemoryCache()
	endpoint := "http://127.0.0.1:1/v3/"
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		TokenCache:           cache,
		CacheNamespace:       "profile",
		RedirectPort:         findAvailablePort(t),
	}
	cacheKey := websso.CacheKey(endpoint, opts)
	cached := fmt.Sprintf(`{"token_id":"cached-token","expires_at":%q,"endpoint":%q}`, time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano), endpoint)
	th.AssertNoErr(t, cache.Set(cacheKey, cached))

	opened := false
	opts.BrowserOpener = func(string) error {
		opened = true
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := websso.Authenticate(ctx, &gophercloud.ServiceClient{
		ProviderClient: newTestProviderClient(),
		Endpoint:       endpoint,
	}, opts)

	if !errors.Is(result.Err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", result.Err)
	}
	if opened {
		t.Fatal("browser opened after cached token validation was canceled")
	}
	if value, err := cache.Get(cacheKey); err != nil || value == "" {
		t.Fatalf("cached token removed after cancellation: value %q, error %v", value, err)
	}
}

func tokenResponse(projectID string) string {
	project := ""
	if projectID != "" {
		project = fmt.Sprintf(`,"project":{"id":%q,"name":%q,"domain":{"id":"default","name":"Default"}}`, projectID, projectID)
	}
	return fmt.Sprintf(`{"token":{"methods":["mapped"],"expires_at":"2035-06-03T02:19:49Z","user":{"id":"user-id","name":"user","domain":{"id":"default","name":"Default"}}%s}}`, project)
}

func newTestProviderClient() *gophercloud.ProviderClient {
	client := new(gophercloud.ProviderClient)
	client.UseTokenLock()
	return client
}

func TestWebSSOCacheHitSkipsBrowser(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()

	endpoint := fakeKeystone.Endpoint()
	tokenID := "websso-cached-token-xyz"
	futureExpiry := time.Now().Add(2 * time.Hour).UTC().Format("2006-01-02T15:04:05.000000Z")

	fakeKeystone.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")

		w.Header().Set("X-Subject-Token", tokenID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"token": {
				"methods": ["token"],
				"expires_at": "%s",
				"catalog": [
					{
						"endpoints": [
							{
								"url": "%scompute/v2.1",
								"interface": "public",
								"region": "RegionOne",
								"region_id": "RegionOne",
								"id": "ep-001"
							}
						],
						"type": "compute",
						"id": "svc-001",
						"name": "nova"
					}
				],
				"user": {
					"id": "user-001",
					"name": "test-user",
					"domain": {"id": "default", "name": "Default"}
				}
			}
		}`, futureExpiry, endpoint)
	})

	cache := newMemoryCache()

	v3Endpoint := endpoint + "v3/"
	client := gophercloud.ServiceClient{
		ProviderClient: newTestProviderClient(),
		Endpoint:       v3Endpoint,
	}

	port := findAvailablePort(t)
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectPort:         port,
		Timeout:              10 * time.Second,
		TokenCache:           cache,
		CacheNamespace:       "test-profile",
	}

	results := startWebSSO(context.Background(), &client, opts)
	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/auth/websso/", port)
	resp, err := http.PostForm(callbackURL, url.Values{"token": {tokenID}})
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.CheckEquals(t, http.StatusOK, resp.StatusCode)

	result1 := <-results
	th.AssertNoErr(t, result1.Err)

	extractedID1, err := result1.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, tokenID, extractedID1)

	cacheKey := websso.CacheKey(v3Endpoint, opts)
	val, err := cache.Get(cacheKey)
	th.AssertNoErr(t, err)
	if val == "" {
		t.Fatal("Expected token to be cached after first WebSSO call")
	}

	opts2 := *opts
	opts2.RedirectPort = findAvailablePort(t)

	result2 := websso.Authenticate(context.Background(), &client, &opts2)
	th.AssertNoErr(t, result2.Err)

	extractedID2, err := result2.ExtractTokenID()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, tokenID, extractedID2)

	catalog, err := result2.ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	if catalog == nil || len(catalog.Entries) == 0 {
		t.Fatal("Expected service catalog from cached token")
	}
	th.AssertEquals(t, "compute", catalog.Entries[0].Type)
}
