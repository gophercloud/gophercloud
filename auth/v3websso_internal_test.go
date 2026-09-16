package auth

import (
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestWebSSOEndpointOrigin(t *testing.T) {
	tests := map[string]string{
		"https://Keystone.Example/identity/v3": "https://keystone.example",
		"https://keystone.example:443/v3":      "https://keystone.example",
		"http://keystone.example:80/v3":        "http://keystone.example",
		"https://keystone.example:5000/v3":     "https://keystone.example:5000",
		"http://keystone.example:443/v3":       "http://keystone.example:443",
		"https://[::1]:443/v3":                 "https://[::1]",
		"https://[::1]:5000/v3":                "https://[::1]:5000",
	}
	for rawURL, expected := range tests {
		actual, err := endpointOrigin(rawURL)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, expected, actual)
	}

	_, err := endpointOrigin("keystone.example/v3")
	th.AssertErr(t, err)
}

// The callback cannot listen on these ports in a test.
func TestV3WebSSOOptsAcceptsLoopbackAddresses(t *testing.T) {
	for _, opts := range []V3WebSSOOpts{
		{},
		{RedirectHost: "localhost", RedirectPort: 1},
		{RedirectHost: "127.0.0.1", RedirectPort: 65535},
		{RedirectHost: "::1", RedirectPort: 9990},
	} {
		opts.IdentityProviderName = "my-idp"
		opts.Protocol = "openid"
		th.AssertNoErr(t, opts.validate())
	}
}
