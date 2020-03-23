package services

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-services")
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("os-services", id)
}
