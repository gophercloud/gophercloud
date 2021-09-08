package speaker

import "github.com/gophercloud/gophercloud"

const baseurl = "bgp-speakers"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseurl, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseurl)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
