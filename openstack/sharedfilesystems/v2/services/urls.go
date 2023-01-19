package services

import "github.com/bizflycloud/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("services")
}
