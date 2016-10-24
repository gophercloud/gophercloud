package shares

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("shares")
}
