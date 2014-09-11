package containers

import "github.com/rackspace/gophercloud"

// accountURL returns the URI used to list Containers.
func accountURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

// containerURL returns the URI for making Container requests.
func containerURL(c *gophercloud.ServiceClient, container string) string {
	return c.ServiceURL(container)
}
