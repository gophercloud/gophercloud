package speaker

import "github.com/gophercloud/gophercloud"

const urlBase = "bgp-speakers"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(urlBase, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(urlBase)
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
	return c.ServiceURL(urlBase, id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(urlBase, id)
}
