package clustertemplates

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("clustertemplates")
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("clustertemplates")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clustertemplates", id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clustertemplates", id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clustertemplates", id)
}
