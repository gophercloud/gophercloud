package v2

import (
	"fmt"

	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

// Client is a structure that contains information for communicating with a provider.
type Client struct {
	endpoint  string
	authority identity.AuthResults
	options   identity.AuthOptions
	token     *identity.Token
}

// NewClient creates and returns a *Client.
func NewClient(e string, a identity.AuthResults, o identity.AuthOptions) *Client {
	return &Client{
		endpoint:  e,
		authority: a,
		options:   o,
	}
}

// ListNetworksURL returns the URL for listing networks available to the tenant.
func (c *Client) ListNetworksURL() string {
	return fmt.Sprintf("%sv2.0/networks", c.endpoint)
}

// CreateNetworkURL returns the URL for creating a network.
func (c *Client) CreateNetworkURL() string {
	return c.ListNetworksURL()
}

// GetNetworkURL returns the URL for showing information for the network with the given id.
func (c *Client) GetNetworkURL(id string) string {
	return fmt.Sprintf("%sv2.0/networks/%s", c.endpoint, id)
}

// UpdateNetworkURL returns the URL for updating information for the network with the given id.
func (c *Client) UpdateNetworkURL(id string) string {
	return c.GetNetworkURL(id)
}

// DeleteNetworkURL returns the URL for deleting the network with the given id.
func (c *Client) DeleteNetworkURL(id string) string {
	return c.GetNetworkURL(id)
}

func (c *Client) ListSubnetsURL() string {
	return fmt.Sprintf("%sv2.0/subnets", c.endpoint)
}

func (c *Client) CreateSubnetURL() string {
	return c.ListSubnetsURL()
}

func (c *Client) DeleteSubnetURL(id string) string {
	return fmt.Sprintf("%sv2.0/subnets/%s", c.endpoint, id)
}

func (c *Client) GetSubnetURL(id string) string {
	return c.DeleteSubnetURL(id)
}

func (c *Client) UpdateSubnetURL(id string) string {
	return c.DeleteSubnetURL(id)
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

	return c.token.Id, err
}
