package groups

import "github.com/rackspace/gophercloud"

const (
	version  = "v2.0"
	rootPath = "security-groups"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, rootPath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, rootPath, id)
}
