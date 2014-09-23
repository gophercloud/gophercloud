package networks

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "networks", id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "networks")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
