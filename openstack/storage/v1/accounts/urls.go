package accounts

import "github.com/rackspace/gophercloud"

// accountURL returns the URI for making Account requests.
func accountURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}
