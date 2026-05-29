package swauth

import "github.com/gophercloud/gophercloud/v2"

func getURL(c *gophercloud.ProviderClient) string {
	return c.IdentityBase + "auth/v1.0"
}
