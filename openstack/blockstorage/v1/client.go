package blockstorage

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/identity"
)

// Client abstracts the connection information needed to make API requests for OpenStack compute endpoints.
type Client struct {
	endpoint  string
	authority identity.AuthResults
	options   identity.AuthOptions
	token     *identity.Token
}

// NewClient creates a new Client structure to use when issuing requests to the server.
func NewClient(e string, a identity.AuthResults, o identity.AuthOptions) *Client {
	return &Client{
		endpoint:  e,
		authority: a,
		options:   o,
	}
}

func (c *Client) GetVolumeURL() string {
	return fmt.Sprintf("%s/volumes", c.endpoint)
}

func (c *Client) GetHeaders() (map[string]string, error) {
	t, err := c.getAuthToken()
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"X-Auth-Token": t,
	}, nil
}

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
