package pools

import "github.com/rackspace/gophercloud"

const (
	version      = "v2.0"
	rootPath     = "lb"
	resourcePath = "pools"
	monitorPath  = "health_monitors"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, rootPath, resourcePath, id)
}

func associateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, rootPath, resourcePath, id, monitorPath)
}

func disassociateURL(c *gophercloud.ServiceClient, poolID, monitorID string) string {
	return c.ServiceURL(version, rootPath, resourcePath, poolID, monitorPath, monitorID)
}
