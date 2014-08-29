package tokens

import "github.com/rackspace/gophercloud"

func getTokenURL(c *gophercloud.ProviderClient) string {
	return c.ServiceURL("auth", "tokens")
}
