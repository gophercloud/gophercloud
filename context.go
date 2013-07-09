package gophercloud

import (
	"net/http"
)

// Provider structures exist for each tangible provider of OpenStack service.
// For example, Rackspace, Hewlett-Packard, and NASA might have their own instance of this structure.
//
// At a minimum, a provider must expose an authentication endpoint.
type Provider struct {
	AuthEndpoint string
}

// Context structures encapsulate Gophercloud-global state in a manner which
// facilitates easier unit testing.  As a user of this SDK, you'll never
// have to use this structure, except when contributing new code to the SDK.
type Context struct {
	// providerMap serves as a directory of supported providers.
	providerMap map[string]Provider

	// httpClient refers to the current HTTP client interface to use.
	httpClient *http.Client
}

// TestContext yields a new Context instance, pre-initialized with a barren
// state suitable for per-unit-test customization.  This configuration consists
// of:
//
// * An empty provider map.
//
// * An HTTP client built by the net/http package (see http://godoc.org/net/http#Client).
func TestContext() *Context {
	return &Context{
		providerMap: make(map[string]Provider),
		httpClient:  &http.Client{},
	}
}

// UseCustomClient configures the context to use a customized HTTP client
// instance.  By default, TestContext() will return a Context which uses
// the net/http package's default client instance.
func (c *Context) UseCustomClient(hc *http.Client) *Context {
	c.httpClient = hc
	return c
}

// RegisterProvider allows a unit test to register a mythical provider convenient for testing.
// If the provider structure lacks adequate configuration, or the configuration given has some
// detectable error, an ErrConfiguration error will result.
func (c *Context) RegisterProvider(name string, p Provider) error {
	if p.AuthEndpoint == "" {
		return ErrConfiguration
	}

	c.providerMap[name] = p
	return nil
}

// WithProvider offers convenience for unit tests.
func (c *Context) WithProvider(name string, p Provider) *Context {
	err := c.RegisterProvider(name, p)
	if err != nil {
		panic(err)
	}
	return c
}

// ProviderByName will locate a provider amongst those previously registered, if it exists.
// If the named provider has not been registered, an ErrProvider error will result.
func (c *Context) ProviderByName(name string) (p Provider, err error) {
	for provider, descriptor := range c.providerMap {
		if name == provider {
			return descriptor, nil
		}
	}
	return Provider{}, ErrProvider
}

// Instantiates a Cloud Servers API for the provider given.
func (c *Context) ServersApi(acc AccessProvider, criteria ApiCriteria) (CloudServersProvider, error) {
	url := acc.FirstEndpointUrlByCriteria(criteria)
	if url == "" {
		return nil, ErrEndpoint
	}

	gcp := &genericServersProvider{
		endpoint: url,
		context:  c,
		access:   acc,
	}

	return gcp, nil
}
