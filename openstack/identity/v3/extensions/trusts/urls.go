package trusts

import "github.com/gophercloud/gophercloud"

const resourcePath = "OS-TRUST/trusts"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
