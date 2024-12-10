package remoteconsoles

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath = "servers"

	resourcePath = "remote-consoles"
)

func rootURL(c gophercloud.Client, serverID string) string {
	return c.ServiceURL(rootPath, serverID, resourcePath)
}

func createURL(c gophercloud.Client, serverID string) string {
	return rootURL(c, serverID)
}
