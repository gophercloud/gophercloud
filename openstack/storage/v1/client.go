package v1

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

// Client is a structure that contains information for communicating with a provider.
type Client struct {
	endpoint  string
	authority identity.AuthResults
	options   gophercloud.AuthOptions
	token     *identity.Token
}

// NewClient creates and returns a *Client.
func NewClient(e string, a identity.AuthResults, o gophercloud.AuthOptions) *Client {
	return &Client{
		endpoint:  e,
		authority: a,
		options:   o,
	}
}

// GetAccountURL returns the URI for making Account requests. This function is exported to allow
// the 'Accounts' subpackage to use it. It is not meant for public consumption.
func (c *Client) GetAccountURL() string {
	return fmt.Sprintf("%s", c.endpoint)
}

// GetContainerURL returns the URI for making Container requests. This function is exported to allow
// the 'Containers' subpackage to use it. It is not meant for public consumption.
func (c *Client) GetContainerURL(container string) string {
	return fmt.Sprintf("%s/%s", c.endpoint, container)
}

// GetObjectURL returns the URI for making Object requests. This function is exported to allow
// the 'Objects' subpackage to use it. It is not meant for public consumption.
func (c *Client) GetObjectURL(container, object string) string {
	return fmt.Sprintf("%s/%s/%s", c.endpoint, container, object)
}

// GetHeaders is a function that gets the header for token authentication against a client's endpoint.
// This function is exported to allow the subpackages to use it. It is not meant for public consumption.
func (c *Client) GetHeaders() (map[string]string, error) {
	t, err := c.getAuthToken()
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"X-Auth-Token": t,
	}, nil
}

// getAuthToken is a function that tries to retrieve an authentication token from a client's endpoint.
func (c *Client) getAuthToken() (string, error) {
	var err error

	if c.token == nil {
		c.token, err = identity.GetToken(c.authority)
		if err != nil {
			return "", err
		}
	}

	return c.token.ID, err
}
