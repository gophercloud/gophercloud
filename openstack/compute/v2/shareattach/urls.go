package shareattach

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "shares"

func resourceURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, resourcePath)
}

func listURL(c *gophercloud.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func createURL(c *gophercloud.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func getURL(c *gophercloud.ServiceClient, serverID, shareID string) string {
	return c.ServiceURL("servers", serverID, resourcePath, shareID)
}

func deleteURL(c *gophercloud.ServiceClient, serverID, shareID string) string {
	return getURL(c, serverID, shareID)
}
