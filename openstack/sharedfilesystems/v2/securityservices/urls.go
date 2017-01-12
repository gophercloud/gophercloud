package securityservices

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("security-services")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("security-services", id)
}
