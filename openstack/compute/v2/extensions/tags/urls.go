package tags

import "github.com/gophercloud/gophercloud"

const (
	rootResourcePath = "servers"
	resourcePath     = "tags"
)

func rootURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL(rootResourcePath, serverID, resourcePath)
}

func listURL(c *gophercloud.ServiceClient, serverID string) string {
	return rootURL(c, serverID)
}
