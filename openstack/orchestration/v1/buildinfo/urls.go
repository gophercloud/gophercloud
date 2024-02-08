package buildinfo

import "github.com/gophercloud/gophercloud/v2"

func getURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("build_info")
}
