package flavors

import (
	"github.com/rackspace/gophercloud/openstack/identity"
)

type Client struct {
	endpoint string
	authority identity.AuthResults
	options identity.AuthOptions
}

func NewClient(e string, a identity.AuthResults, ao identity.AuthOptions) *Client {
	return &Client{
		endpoint: e,
		authority: a,
		options: ao,
	}
}

func (c *Client) getListUrl() string {
	return c.endpoint + "/flavors/detail"
}

func (c *Client) getListHeaders() (map[string]string, error) {
	t, err := identity.GetToken(c.authority)
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"X-Auth-Token": t.Id,
	}, nil
}

