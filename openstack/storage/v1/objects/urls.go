package objects

import "github.com/rackspace/gophercloud"

// getObjectURL returns the URI for making Object requests.
func getObjectURL(c *gophercloud.ServiceClient, container, object string) string {
	return c.ServiceURL(container, object)
}

// getContainerURL returns the URI for making Container requests.
func getContainerURL(c *gophercloud.ServiceClient, container string) string {
	return c.ServiceURL(container)
}
