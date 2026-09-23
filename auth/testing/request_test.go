package testing

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// v3AuthRequest exercises the public authentication path without a network call.
// The fixed URL keeps OAuth signatures deterministic.
func v3AuthRequest(t *testing.T, opts auth.AuthOptionsBuilderV3) (map[string]any, http.Header, error) {
	t.Helper()
	var body map[string]any
	var headers http.Header
	provider := &gophercloud.ProviderClient{HTTPClient: http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			headers = req.Header.Clone()
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusCreated,
				Header:     http.Header{"X-Subject-Token": []string{"testtoken"}},
				Body:       io.NopCloser(strings.NewReader(`{"token": {}}`)),
			}, nil
		}),
	}}
	_, err := (auth.AuthOptionsV3{AuthURL: "http://127.0.0.1:33199/v3", Auth: opts}).Authenticate(context.Background(), provider)
	return body, headers, err
}
