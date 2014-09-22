package extensions

import "github.com/rackspace/gophercloud"

const version = "v2.0"

func extensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL(version, "extensions", name)
}

func listExtensionURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, "extensions")
}
