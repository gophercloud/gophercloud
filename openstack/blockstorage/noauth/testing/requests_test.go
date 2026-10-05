package testing

import (
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/noauth"
)

type customNoAuthOpts string

func (opts customNoAuthOpts) ToNoAuthToken() (string, error) {
	return string(opts), nil
}

func TestNoAuthClient(t *testing.T) {
	tests := []struct {
		name string
		opts auth.NoAuthOpts
		want string
	}{
		{
			name: "v2 explicit",
			opts: auth.NoAuthV2Opts{Username: "user", TenantName: "tenant"},
			want: "user:tenant",
		},
		{
			name: "v2 defaults",
			opts: auth.NoAuthV2Opts{},
			want: "admin:admin",
		},
		{
			name: "v3 explicit",
			opts: auth.NoAuthV3Opts{UserID: "user-id", ProjectID: "project-id"},
			want: "user-id:project-id",
		},
		{
			name: "v3 project defaults to user",
			opts: auth.NoAuthV3Opts{UserID: "user-id"},
			want: "user-id:user-id",
		},
		{
			name: "custom options",
			opts: customNoAuthOpts("custom:scope"),
			want: "custom:scope",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider, err := noauth.NewClient(test.opts)
			if err != nil {
				t.Fatal(err)
			}
			if provider.TokenID != test.want {
				t.Fatalf("token = %q, want %q", provider.TokenID, test.want)
			}
		})
	}
}

func TestNoAuthMalformedToken(t *testing.T) {
	for _, token := range []string{"", "user", "user:", ":project", "user:project:extra"} {
		t.Run(token, func(t *testing.T) {
			provider, err := noauth.NewClient(customNoAuthOpts(token))
			if err != nil {
				t.Fatal(err)
			}
			_, err = noauth.NewBlockStorageNoAuthV3(provider, noauth.EndpointOpts{CinderEndpoint: "http://cinder:8776/v3"})
			if err == nil {
				t.Fatal("expected malformed token error")
			}
			if token != "" && strings.Contains(err.Error(), token) {
				t.Fatalf("error discloses token %q", token)
			}
		})
	}
}

func TestNoAuthV2Endpoint(t *testing.T) {
	tests := []struct {
		name     string
		opts     auth.NoAuthV2Opts
		endpoint string
		want     NoAuthResult
	}{
		{
			name:     "explicit",
			opts:     auth.NoAuthV2Opts{Username: "user", TenantName: "test"},
			endpoint: "http://cinder:8776/v2",
			want:     naTestResult,
		},
		{
			name:     "defaults",
			opts:     auth.NoAuthV2Opts{},
			endpoint: "http://cinder:8776/v2/",
			want:     naResult,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider, err := noauth.NewClient(test.opts)
			if err != nil {
				t.Fatal(err)
			}
			client, err := noauth.NewBlockStorageNoAuthV2(provider, noauth.EndpointOpts{CinderEndpoint: test.endpoint})
			if err != nil {
				t.Fatal(err)
			}
			if client.Endpoint != test.want.Endpoint {
				t.Fatalf("endpoint = %q, want %q", client.Endpoint, test.want.Endpoint)
			}
			if client.TokenID != test.want.TokenID {
				t.Fatalf("token = %q, want %q", client.TokenID, test.want.TokenID)
			}
			if client.Type != "block-storage" {
				t.Fatalf("type = %q", client.Type)
			}
		})
	}
}

func TestNoAuthV3Endpoint(t *testing.T) {
	provider, err := noauth.NewClient(auth.NoAuthV3Opts{UserID: "user-id", ProjectID: "project-id"})
	if err != nil {
		t.Fatal(err)
	}
	client, err := noauth.NewBlockStorageNoAuthV3(provider, noauth.EndpointOpts{CinderEndpoint: "http://cinder:8776/v3"})
	if err != nil {
		t.Fatal(err)
	}
	if client.Endpoint != "http://cinder:8776/v3/project-id/" {
		t.Fatalf("endpoint = %q", client.Endpoint)
	}
	if client.TokenID != "user-id:project-id" {
		t.Fatalf("token = %q", client.TokenID)
	}
}

func TestNoAuthV3RequiresEndpoint(t *testing.T) {
	provider, err := noauth.NewClient(auth.NoAuthV3Opts{UserID: "user-id"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = noauth.NewBlockStorageNoAuthV3(provider, noauth.EndpointOpts{})
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != errorResult {
		t.Fatalf("error = %q, want %q", err, errorResult)
	}
}
