package bootfromvolume

import "github.com/gophercloud/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("servers")
}
