package tools

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CheckSinglePage checks that EachPage and AllPages retrieve the exact expected
// IDs with one collection GET each. The client must have a dedicated provider
// and must not be used concurrently while this helper replaces its transport.
func CheckSinglePage(t *testing.T, ctx context.Context, client *gophercloud.ServiceClient, collectionURL string, newPager func() pagination.Pager, extractIDs func(pagination.Page) ([]string, error), expectedIDs []string) {
	t.Helper()
	seen := make(map[string]bool)
	for _, id := range expectedIDs {
		if id == "" || seen[id] {
			t.Fatalf("expected IDs must be distinct and nonempty: %v", expectedIDs)
		}
		seen[id] = true
	}
	u, err := url.Parse(collectionURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.RawQuery != "" {
		t.Fatalf("invalid collection URL %q", collectionURL)
	}
	for _, allPages := range []bool{false, true} {
		name := "EachPage"
		if allPages {
			name = "AllPages"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()
			original := client.HTTPClient.Transport
			next := original
			if next == nil {
				next = http.DefaultTransport
			}
			rt := &singlePageTransport{collection: *u, next: next}
			client.HTTPClient.Transport = rt
			t.Cleanup(func() { client.HTTPClient.Transport = original })

			var actual []string
			collect := func(page pagination.Page) error {
				ids, err := extractIDs(page)
				actual = append(actual, ids...)
				return err
			}
			pager := newPager()
			if allPages {
				page, err := pager.AllPages(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if err := collect(page); err != nil {
					t.Fatal(err)
				}
			} else {
				pages := 0
				err := pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					if pages > 1 {
						return false, fmt.Errorf("expected a single nonempty page")
					}
					return true, collect(page)
				})
				if err != nil {
					t.Fatal(err)
				}
			}
			expected := slices.Clone(expectedIDs)
			slices.Sort(expected)
			slices.Sort(actual)
			if !slices.Equal(expected, actual) {
				t.Errorf("expected IDs %v, got %v", expected, actual)
			}
			if got := rt.requests.Load(); got != 1 {
				t.Errorf("expected one collection GET, got %d", got)
			}
		})
	}
}

type singlePageTransport struct {
	// collection identifies the collection whose requests are counted.
	collection url.URL
	// next sends requests without modifying them.
	next http.RoundTripper
	// requests counts matching GETs, including rejected continuations.
	requests atomic.Int64
}

func (rt *singlePageTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodGet && req.URL.Scheme == rt.collection.Scheme && req.URL.Host == rt.collection.Host && req.URL.EscapedPath() == rt.collection.EscapedPath() {
		if rt.requests.Add(1) > 1 {
			return nil, fmt.Errorf("unexpected continuation request for a single-page collection")
		}
	}
	return rt.next.RoundTrip(req)
}
