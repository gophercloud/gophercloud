package testing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGenericAuthDiscoveryUsesExactEndpoint(t *testing.T) {
	for _, source := range []string{"cloud", "environment"} {
		t.Run(source, func(t *testing.T) {
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/":
					fmt.Fprintf(w, `{"versions":{"values":[{"id":"v3.0","status":"stable","links":[{"rel":"self","href":%q}]}]}}`, server.URL+"/nonstandard-auth/")
				case "/nonstandard-auth/auth/tokens":
					w.Header().Set("X-Subject-Token", "discovered-token")
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, `{"token":{}}`)
				default:
					t.Errorf("unexpected path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			var plugin auth.Authenticator
			var err error
			switch source {
			case "cloud":
				plugin, err = auth.AuthOptionsFromCloud(cloudSource{
					authType: auth.AuthPassword,
					authData: map[string]any{
						"auth_url": server.URL,
						"user_id":  "user",
						"password": "password",
					},
				})
			case "environment":
				t.Setenv("OS_AUTH_TYPE", "password")
				t.Setenv("OS_AUTH_URL", server.URL)
				t.Setenv("OS_USERID", "user")
				t.Setenv("OS_PASSWORD", "password")
				plugin, err = auth.AuthOptionsFromEnv()
			}
			th.AssertNoErr(t, err)

			result, err := plugin.Authenticate(context.Background(), nil)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, "discovered-token", result.TokenID)
		})
	}
}

func TestGenericAuthDiscoveryUsesProviderRetryPolicy(t *testing.T) {
	var discoveryCalls atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/":
			if discoveryCalls.Add(1) == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, `{"versions":{"values":[{"id":"v3.0","status":"stable","links":[{"rel":"self","href":%q}]}]}}`, server.URL+"/v3/")
		case "/v3/auth/tokens":
			w.Header().Set("X-Subject-Token", "discovered-token")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"token":{}}`)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Setenv("OS_AUTH_TYPE", "password")
	t.Setenv("OS_AUTH_URL", server.URL)
	t.Setenv("OS_USERID", "user")
	t.Setenv("OS_PASSWORD", "password")
	plugin, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	var retryCalls atomic.Int32
	provider := &gophercloud.ProviderClient{
		MaxBackoffRetries: 2,
		RetryFunc: func(context.Context, string, string, *gophercloud.RequestOpts, error, uint) error {
			retryCalls.Add(1)
			return nil
		},
	}
	result, err := plugin.Authenticate(context.Background(), provider)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "discovered-token", result.TokenID)
	th.AssertEquals(t, int32(2), discoveryCalls.Load())
	th.AssertEquals(t, int32(1), retryCalls.Load())
}

func TestGenericAuthDiscoveryUsesProviderAndContext(t *testing.T) {
	for _, source := range []string{"cloud", "environment"} {
		t.Run(source, func(t *testing.T) {
			var calls atomic.Int32
			var server *httptest.Server
			server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/":
					fmt.Fprintf(w, `{"versions":{"values":[{"id":"v3.0","status":"stable","links":[{"rel":"self","href":%q}]}]}}`, server.URL+"/identity/v3/")
				case "/identity/v3/auth/tokens":
					w.Header().Set("X-Subject-Token", "discovered-token")
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, `{"token":{"user":{"id":"user"},"expires_at":"2099-01-01T00:00:00Z"}}`)
				default:
					t.Errorf("unexpected path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			var plugin auth.Authenticator
			var err error
			switch source {
			case "cloud":
				plugin, err = auth.AuthOptionsFromCloud(cloudSource{
					authType: auth.AuthPassword,
					authData: map[string]any{
						"auth_url": server.URL,
						"user_id":  "user",
						"password": "password",
					},
				})
			case "environment":
				t.Setenv("OS_AUTH_TYPE", "password")
				t.Setenv("OS_AUTH_URL", server.URL)
				t.Setenv("OS_USERID", "user")
				t.Setenv("OS_PASSWORD", "password")
				plugin, err = auth.AuthOptionsFromEnv()
			}
			if err != nil {
				t.Fatal(err)
			}
			if calls.Load() != 0 {
				t.Fatal("auth option construction performed network I/O")
			}

			provider := &gophercloud.ProviderClient{HTTPClient: *server.Client()}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if _, err = plugin.Authenticate(ctx, provider); !errors.Is(err, context.Canceled) {
				t.Fatalf("expected cancellation, got %v", err)
			}

			result, err := plugin.Authenticate(context.Background(), provider)
			if err != nil {
				t.Fatal(err)
			}
			if result.TokenID != "discovered-token" {
				t.Fatalf("expected discovered-token, got %q", result.TokenID)
			}
			if calls.Load() != 2 {
				t.Fatalf("expected two requests, got %d", calls.Load())
			}
		})
	}
}
