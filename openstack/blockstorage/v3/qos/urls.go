package qos

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("qos-specs")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("qos-specs")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("qos-specs", id)
}
