package monitors

import "github.com/rackspace/gophercloud"

const (
	version      = "v2.0"
	rootPath     = "lb"
	resourcePath = "healthmonitors"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, rootPath, resourcePath, id)
}
