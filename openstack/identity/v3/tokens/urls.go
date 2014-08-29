package tokens

import "github.com/rackspace/gophercloud"

func getTokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
