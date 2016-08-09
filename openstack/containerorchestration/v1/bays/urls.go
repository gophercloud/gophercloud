package bays

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("bays")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("bays")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("bays", id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("bays", id)
}
