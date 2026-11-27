package tools

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestListLimitTransport(t *testing.T) {
	collection, err := url.Parse("https://example.com/v2/items")
	th.AssertNoErr(t, err)
	for _, tc := range []struct {
		name   string
		method string
		url    string
		query  string
		err    string
	}{
		{name: "initial", method: "GET", url: collection.String(), query: "limit=1"},
		{name: "preserve filters", method: "GET", url: collection.String() + "?name=test&sort_key=id", query: "limit=1&name=test&sort_key=id"},
		{name: "preserve marker", method: "GET", url: collection.String() + "?limit=1&marker=abc", query: "limit=1&marker=abc"},
		{name: "preserve offset", method: "GET", url: collection.String() + "?limit=1&offset=2", query: "limit=1&offset=2"},
		{name: "lost marker limit", method: "GET", url: collection.String() + "?marker=abc", err: "lost limit=1"},
		{name: "lost offset limit", method: "GET", url: collection.String() + "?offset=2", err: "lost limit=1"},
		{name: "lost reverse limit", method: "GET", url: collection.String() + "?page_reverse=true", err: "lost limit=1"},
		{name: "changed limit", method: "GET", url: collection.String() + "?limit=2", err: "must preserve limit=1"},
		{name: "empty limit", method: "GET", url: collection.String() + "?limit=", err: "must preserve limit=1"},
		{name: "multiple limits", method: "GET", url: collection.String() + "?limit=1&limit=2", err: "must preserve limit=1"},
		{name: "other method", method: "POST", url: collection.String()},
		{name: "other path", method: "GET", url: collection.String() + "/abc"},
		{name: "other host", method: "GET", url: "https://other.example.com/v2/items"},
		{name: "other scheme", method: "GET", url: "http://example.com/v2/items"},
		{name: "trailing slash", method: "GET", url: collection.String() + "/"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), tc.method, tc.url, nil)
			th.AssertNoErr(t, err)
			req.Header.Set("X-Test", "preserved")
			originalURL := req.URL.String()
			called := false
			rt := listLimitTransport{collection: *collection, next: roundTripFunc(func(actual *http.Request) (*http.Response, error) {
				called = true
				th.CheckEquals(t, tc.method, actual.Method)
				th.CheckEquals(t, tc.query, actual.URL.RawQuery)
				th.CheckEquals(t, "preserved", actual.Header.Get("X-Test"))
				th.CheckEquals(t, req.Context(), actual.Context())
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}")), Request: actual}, nil
			})}
			resp, err := rt.RoundTrip(req)
			if tc.err != "" {
				th.AssertErr(t, err)
				th.CheckTrue(t, strings.Contains(err.Error(), tc.err))
				th.CheckFalse(t, called)
			} else {
				th.AssertNoErr(t, err)
				resp.Body.Close()
				th.CheckTrue(t, called)
			}
			th.CheckEquals(t, originalURL, req.URL.String())
		})
	}
}

func TestLimitListPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.CheckEquals(t, "1", r.URL.Query().Get("limit"))
		th.CheckEquals(t, "test", r.URL.Query().Get("name"))
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"items": []}`)
	})
	for _, custom := range []bool{false, true} {
		client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: fakeServer.Endpoint()}
		var original http.RoundTripper
		if custom {
			original = &http.Transport{}
			defer original.(*http.Transport).CloseIdleConnections()
			client.HTTPClient.Transport = original
		}
		t.Run("wrapped", func(t *testing.T) {
			LimitListPages(t, client, client.ServiceURL("items"))
			resp, err := client.Get(context.Background(), client.ServiceURL("items")+"?name=test", nil, &gophercloud.RequestOpts{KeepResponseBody: true})
			th.AssertNoErr(t, err)
			page, err := pagination.PageResultFrom(resp)
			th.AssertNoErr(t, err)
			// Offset-based page types must see the injected limit in the response URL.
			th.CheckEquals(t, "1", page.URL.Query().Get("limit"))
			th.CheckEquals(t, "test", page.URL.Query().Get("name"))
		})
		th.AssertEquals(t, original, client.HTTPClient.Transport)
	}
}
