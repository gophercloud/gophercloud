package floatingips

import "github.com/rackspace/gophercloud"

const (
	version      = "v2.0"
	resourcePath = "floatingips"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id)
}
