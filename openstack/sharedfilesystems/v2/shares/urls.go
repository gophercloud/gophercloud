package shares

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("shares")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func listURL(c *gophercloud.ServiceClient, detail bool) string {
	if detail {
		return c.ServiceURL("shares", "detail")
	}
	return c.ServiceURL("shares")
}
