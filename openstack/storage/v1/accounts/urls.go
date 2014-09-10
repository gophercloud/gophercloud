package accounts

import "github.com/rackspace/gophercloud"

// getAccountURL returns the URI for making Account requests.
func getAccountURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}
