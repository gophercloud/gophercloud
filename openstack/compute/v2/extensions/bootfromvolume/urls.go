package bootfromvolume

import "github.com/bizflycloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("servers")
}
