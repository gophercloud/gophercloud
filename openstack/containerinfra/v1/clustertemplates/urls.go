package clustertemplates

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("clustertemplates")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clustertemplates", id)
}
