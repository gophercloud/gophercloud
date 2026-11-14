package testing

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsValidation(t *testing.T) {
	tests := []struct {
		name     string
		argument string
		missing  bool
		modify   func(*websso.AuthOptions)
	}{
		{"missing identity provider", "IdentityProviderName", true, func(opts *websso.AuthOptions) { opts.IdentityProviderName = "" }},
		{"missing protocol", "Protocol", true, func(opts *websso.AuthOptions) { opts.Protocol = "" }},
		{"negative port", "RedirectPort", false, func(opts *websso.AuthOptions) { opts.RedirectPort = -1 }},
		{"port too large", "RedirectPort", false, func(opts *websso.AuthOptions) { opts.RedirectPort = 65536 }},
		{"negative timeout", "Timeout", false, func(opts *websso.AuthOptions) { opts.Timeout = -time.Second }},
		{"cache without namespace", "CacheNamespace", false, func(opts *websso.AuthOptions) { opts.TokenCache = newMemoryCache() }},
		{"non-loopback host", "RedirectHost", false, func(opts *websso.AuthOptions) { opts.RedirectHost = "0.0.0.0" }},
		{"invalid host", "RedirectHost", false, func(opts *websso.AuthOptions) { opts.RedirectHost = "invalid" }},
	}
	validators := map[string]func(*websso.AuthOptions) error{
		"ToTokenV3CreateMap": func(opts *websso.AuthOptions) error {
			_, err := opts.ToTokenV3CreateMap(nil)
			return err
		},
		"Authenticate": func(opts *websso.AuthOptions) error {
			return websso.Authenticate(context.Background(), nil, opts).Err
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for name, validate := range validators {
				t.Run(name, func(t *testing.T) {
					opts := &websso.AuthOptions{IdentityProviderName: "my-idp", Protocol: "openid"}
					test.modify(opts)
					err := validate(opts)
					if test.missing {
						var missing gophercloud.ErrMissingInput
						if !errors.As(err, &missing) {
							t.Fatalf("expected ErrMissingInput, got %v", err)
						}
						th.CheckEquals(t, test.argument, missing.Argument)
					} else if err == nil || !strings.Contains(err.Error(), test.argument) {
						t.Fatalf("expected %s validation error, got %v", test.argument, err)
					}
				})
			}
		})
	}
}

func TestToTokenV3CreateMap(t *testing.T) {
	tests := map[string]websso.AuthOptions{
		"defaults":  {},
		"localhost": {RedirectHost: "localhost", RedirectPort: 1},
		"IPv4":      {RedirectHost: "127.0.0.1", RedirectPort: 65535},
		"IPv6":      {RedirectHost: "::1", RedirectPort: 9990},
		"browser and cache hooks": {
			BrowserOpener: func(string) error {
				t.Fatal("validation must not open the browser")
				return nil
			},
			TokenCache:     nonJSONCache{newMemoryCache()},
			CacheNamespace: "profile",
		},
	}
	for name, opts := range tests {
		t.Run(name, func(t *testing.T) {
			opts.IdentityProviderName = "my-idp"
			opts.Protocol = "openid"
			body, err := opts.ToTokenV3CreateMap(nil)
			th.AssertNoErr(t, err)
			if body != nil {
				t.Fatalf("expected nil body, got %v", body)
			}
		})
	}
}

type nonJSONCache struct{ memoryCache }

func (nonJSONCache) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cache implementation must not be serialized")
}
