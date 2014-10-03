package extensions

import "github.com/rackspace/gophercloud"

func extensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL("extensions", name)
}

func listExtensionURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("extensions")
}
