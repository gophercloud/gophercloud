package gophercloud

import (
	"net/http"
)

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
func (c *Context) UseCustomClient(hc *http.Client) {
	c.httpClient = hc
}
