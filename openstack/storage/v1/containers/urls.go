package containers

import "github.com/rackspace/gophercloud"

// getAccountURL returns the URI used to list Containers.
func getAccountURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

// getContainerURL returns the URI for making Container requests.
func getContainerURL(c *gophercloud.ServiceClient, container string) string {
	return c.ServiceURL(container)
}
