package tools

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
)

// LimitListPages supplies limit=1 for a collection whose public List function
// has no options argument. Only initial GET requests to collectionURL are
// modified. Continuation requests must preserve limit=1 themselves.
//
// The client must have a provider dedicated to the test and must not be used
// concurrently. The original transport is restored on test cleanup.
func LimitListPages(t *testing.T, client *gophercloud.ServiceClient, collectionURL string) {
	t.Helper()
	u, err := url.Parse(collectionURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.RawQuery != "" || u.Fragment != "" {
		t.Fatalf("invalid collection URL %q", collectionURL)
	}
	original := client.HTTPClient.Transport
	next := original
	if next == nil {
		next = http.DefaultTransport
	}
	client.HTTPClient.Transport = listLimitTransport{collection: *u, next: next}
	t.Cleanup(func() {
		client.HTTPClient.Transport = original
	})
}

type listLimitTransport struct {
	// collection identifies the exact origin and path to which the limit applies.
	collection url.URL
	// next is the original transport, including any debug logging.
	next http.RoundTripper
}

func (rt listLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet || req.URL.Scheme != rt.collection.Scheme || req.URL.Host != rt.collection.Host || req.URL.EscapedPath() != rt.collection.EscapedPath() {
		return rt.next.RoundTrip(req)
	}
	q := req.URL.Query()
	if q.Has("limit") {
		if values := q["limit"]; len(values) != 1 || values[0] != "1" {
			return nil, fmt.Errorf("collection request must preserve limit=1")
		}
		return rt.next.RoundTrip(req)
	}
	if q.Has("marker") || q.Has("offset") || q.Has("page_reverse") {
		return nil, fmt.Errorf("collection continuation lost limit=1")
	}
	cloned := req.Clone(req.Context())
	q.Set("limit", "1")
	cloned.URL.RawQuery = q.Encode()
	return rt.next.RoundTrip(cloned)
}
