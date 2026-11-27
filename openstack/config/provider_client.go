package config

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

type options struct {
	httpClient http.Client
	tlsConfig  *tls.Config
}

// WithHTTPClient enables passing a custom http.Client to be used in the
// ProviderClient for authentication and for any further call, for example when
// using a ServiceClient derived from this ProviderClient.
func WithHTTPClient(httpClient http.Client) func(*options) {
	return func(o *options) {
		o.httpClient = httpClient
	}
}

// WithTLSConfig replaces the Transport of the default HTTP client (or of the
// HTTP client passed with WithHTTPClient) with a RoundTripper containing the
// given TLS config.
func WithTLSConfig(tlsConfig *tls.Config) func(*options) {
	return func(o *options) {
		o.tlsConfig = tlsConfig
	}
}

// NewProviderClient logs in to an OpenStack cloud found at the identity
// endpoint specified by the options, acquires a token, and returns a Provider
// Client instance that's ready to operate.
//
// The authenticator selects the Identity API version and authentication endpoint.
// Configure any version discovery performed while loading authentication options
// separately, for example with auth.WithTLSConfig when loading a cloud. The
// options passed here configure authentication and subsequent service requests.
func NewProviderClient(ctx context.Context, authOptions auth.Authenticator, opts ...func(*options)) (*gophercloud.ProviderClient, error) {
	var options options
	for _, apply := range opts {
		apply(&options)
	}

	client, err := openstack.NewClient(authOptions.GetAuthURL())
	if err != nil {
		return nil, err
	}

	if options.tlsConfig != nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = options.tlsConfig
		options.httpClient.Transport = transport
	}
	client.HTTPClient = options.httpClient

	err = openstack.Authenticate(ctx, client, authOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}
