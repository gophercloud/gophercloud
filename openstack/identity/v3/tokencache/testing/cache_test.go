package testing

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokencache"
)

type testCache struct {
	value        string
	getErr       error
	setErr       error
	deleteErr    error
	setCalled    bool
	deleteCalled bool
}

func (c *testCache) Get(string) (string, error) {
	return c.value, c.getErr
}

func (c *testCache) Set(_ string, value string) error {
	c.setCalled = true
	c.value = value
	return c.setErr
}

func (c *testCache) Delete(string) error {
	c.deleteCalled = true
	c.value = ""
	return c.deleteErr
}

func TestKeySeparatesFields(t *testing.T) {
	first := tokencache.Key(tokencache.KeyOptions{
		Flow:      "flow\nprincipal=principal",
		Principal: "",
	})
	second := tokencache.Key(tokencache.KeyOptions{
		Flow:      "flow",
		Principal: "principal\nprincipal=",
	})

	if first == second {
		t.Fatalf("different cache identities produced the same key %q", first)
	}
}

func TestPersistAndLoad(t *testing.T) {
	cache := new(testCache)
	tokencache.Persist(cache, "key", "https://identity.example/v3/", "token", time.Now().Add(time.Hour))

	if !cache.setCalled {
		t.Fatal("Persist did not write the token")
	}
	actual, ok := tokencache.Load(cache, "key", "https://identity.example/v3")
	if !ok || actual != "token" {
		t.Fatalf("Load returned (%q, %t), want (%q, true)", actual, ok, "token")
	}
}

func TestEmptyCacheIsMiss(t *testing.T) {
	for _, cache := range []tokencache.Cache{nil, new(testCache)} {
		if token, ok := tokencache.Load(cache, "key", "endpoint"); ok || token != "" {
			t.Fatalf("Load returned (%q, %t), want cache miss", token, ok)
		}
	}
}

func TestPersistSkipsEmptyValues(t *testing.T) {
	cache := new(testCache)
	tokencache.Persist(nil, "key", "endpoint", "token", time.Now().Add(time.Hour))
	tokencache.Persist(cache, "key", "endpoint", "", time.Now().Add(time.Hour))
	if cache.setCalled {
		t.Fatal("Persist stored an empty token")
	}
}

func TestLoadDeletesInvalidEntries(t *testing.T) {
	future := time.Now().Add(time.Hour)
	tests := []struct {
		name     string
		value    string
		endpoint string
	}{
		{name: "malformed", value: "not JSON", endpoint: "https://identity.example/v3"},
		{name: "empty token", value: cacheEntry("", future), endpoint: "https://identity.example/v3"},
		{name: "expired", value: cacheEntry("token", time.Now().Add(-time.Minute)), endpoint: "https://identity.example/v3"},
		{name: "within expiry margin", value: cacheEntry("token", time.Now().Add(4*time.Minute)), endpoint: "https://identity.example/v3"},
		{name: "endpoint mismatch", value: cacheEntry("token", future), endpoint: "https://other.example/v3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := &testCache{value: test.value}
			if token, ok := tokencache.Load(cache, "key", test.endpoint); ok || token != "" {
				t.Fatalf("Load returned (%q, %t), want cache miss", token, ok)
			}
			if !cache.deleteCalled {
				t.Fatal("Load did not delete the invalid entry")
			}
		})
	}
}

func TestCacheErrorsAreBestEffort(t *testing.T) {
	wantErr := errors.New("cache unavailable")
	if token, ok := tokencache.Load(&testCache{getErr: wantErr}, "key", "endpoint"); ok || token != "" {
		t.Fatalf("Load returned (%q, %t) after Get error", token, ok)
	}

	malformed := &testCache{value: "not JSON", deleteErr: wantErr}
	if token, ok := tokencache.Load(malformed, "key", "endpoint"); ok || token != "" {
		t.Fatalf("Load returned (%q, %t) after Delete error", token, ok)
	}
	if !malformed.deleteCalled {
		t.Fatal("Load did not attempt to delete the malformed entry")
	}

	unavailable := &testCache{setErr: wantErr}
	tokencache.Persist(unavailable, "key", "endpoint", "token", time.Now().Add(time.Hour))
	if !unavailable.setCalled {
		t.Fatal("Persist did not attempt to store the token")
	}
}

func cacheEntry(token string, expiresAt time.Time) string {
	return fmt.Sprintf(`{"token_id":%q,"expires_at":%q,"endpoint":%q}`,
		token, expiresAt.UTC().Format(time.RFC3339Nano), "https://identity.example/v3")
}
