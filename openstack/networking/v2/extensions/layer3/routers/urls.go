package routers

import "github.com/rackspace/gophercloud"

const (
	version      = "v2.0"
	resourcePath = "routers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id)
}

func addInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id, "add_router_interface")
}

func removeInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id, "remove_router_interface")
}
