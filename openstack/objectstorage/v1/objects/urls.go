package objects

import "github.com/rackspace/gophercloud"

// objectURL returns the URI for making Object requests.
func objectURL(c *gophercloud.ServiceClient, container, object string) string {
	return c.ServiceURL(container, object)
}

// containerURL returns the URI for making Container requests.
func containerURL(c *gophercloud.ServiceClient, container string) string {
	return c.ServiceURL(container)
}
