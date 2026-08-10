package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokencache"
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

func TestCacheKeyIsNonEmptyAndFlowIsolated(t *testing.T) {
	endpoint := "https://keystone.example.com/v3"
	webSSOKey := websso.CacheKey(endpoint, &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		CacheNamespace:       "profile",
	})
	if webSSOKey == "" {
		t.Fatal("expected non-empty WebSSO cache key")
	}
	oidcKey := tokencache.Key(tokencache.KeyOptions{
		Flow:             "oidc-client-credentials",
		IdentityEndpoint: endpoint,
		Principal:        "profile",
	})
	if webSSOKey == oidcKey {
		t.Fatalf("OIDC and WebSSO must not share cache key %q", webSSOKey)
	}
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
