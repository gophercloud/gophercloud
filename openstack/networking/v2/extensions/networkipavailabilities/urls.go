package networkipavailabilities

import "github.com/gophercloud/gophercloud"

const resourcePath = "network-ip-availabilities"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
