package services

import "github.com/gophercloud/gophercloud/v2"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-services")
}
