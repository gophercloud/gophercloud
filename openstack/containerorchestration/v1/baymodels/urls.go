package baymodels

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("baymodels")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("baymodels", id)
}
