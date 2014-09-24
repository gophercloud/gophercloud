package subnets

import "github.com/rackspace/gophercloud"

const version = "v2.0"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, "subnets", id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, "subnets")
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
