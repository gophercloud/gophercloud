package gophercloud

import (
	"net/http"
)

type Context struct {
	// providerMap serves as a directory of supported providers.
	providerMap map[string]*Provider

	// httpClient refers to the current HTTP client interface to use.
	httpClient *http.Client
}

func TestContext() *Context {
	return &Context{
		providerMap: make(map[string]*Provider),
		httpClient: &http.Client{},
	}
}

func (c *Context) UseCustomClient(hc *http.Client) {
	c.httpClient = hc;
}
