package provider

import "github.com/gophercloud/gophercloud"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("networks")
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
