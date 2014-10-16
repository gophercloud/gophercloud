package cdncontainers

import "github.com/rackspace/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

func enableURL(c *gophercloud.ServiceClient, containerName string) string {
	return c.ServiceURL(containerName)
}

func getURL(c *gophercloud.ServiceClient, containerName string) string {
	return c.ServiceURL(containerName)
}

func updateURL(c *gophercloud.ServiceClient, containerName string) string {
	return c.ServiceURL(containerName)
}
