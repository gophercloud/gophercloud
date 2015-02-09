package instances

import "github.com/rackspace/gophercloud"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("instances")
}

func createURL(c *gophercloud.ServiceClient) string {
	return baseURL(c)
}
