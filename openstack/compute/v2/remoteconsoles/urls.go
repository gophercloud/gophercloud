package remoteconsoles

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath = "servers"

	resourcePath = "remote-consoles"

	consolePath = "os-console-auth-tokens"
)

func rootURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL(rootPath, serverID, resourcePath)
}

func createURL(c *gophercloud.ServiceClient, serverID string) string {
	return rootURL(c, serverID)
}

func getConsoleURL(c *gophercloud.ServiceClient, console string) string {
	return c.ServiceURL(consolePath, console)
}
