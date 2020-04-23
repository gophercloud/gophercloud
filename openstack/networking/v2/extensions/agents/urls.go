package agents

import "github.com/gophercloud/gophercloud"

const resourcePath = "agents"
const dhcpNetworksResourcePath = "dhcp-networks"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func dhcpNetworksURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, dhcpNetworksResourcePath)
}

func listDHCPNetworksURL(c *gophercloud.ServiceClient, id string) string {
	return dhcpNetworksURL(c, id)
}

func scheduleDHCPNetworkURL(c *gophercloud.ServiceClient, id string) string {
	return dhcpNetworksURL(c, id)
}

func removeDHCPNetworkURL(c *gophercloud.ServiceClient, id string, networkID string) string {
	return c.ServiceURL(resourcePath, id, dhcpNetworksResourcePath, networkID)
}
