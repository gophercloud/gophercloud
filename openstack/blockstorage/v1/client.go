package v1

import (
	"fmt"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

// Client abstracts the connection information needed to make API requests for OpenStack block storage endpoints.
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

func (c *Client) GetVolumesURL() string {
	return fmt.Sprintf("%s/volumes", c.endpoint)
}

func (c *Client) GetVolumeURL(id string) string {
	return fmt.Sprintf("%s/volumes/%s", c.endpoint, id)
}

func (c *Client) GetSnapshotsURL() string {
	return fmt.Sprintf("%s/snapshots", c.endpoint)
}

func (c *Client) GetSnapshotURL(id string) string {
	return fmt.Sprintf("%s/snapshots/%s", c.endpoint, id)
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
