package addressscopes

import "github.com/gophercloud/gophercloud"

const resourcePath = "addressscopes"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
