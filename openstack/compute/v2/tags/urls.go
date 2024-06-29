package tags

import "github.com/gophercloud/gophercloud/v2"

const (
	rootResourcePath = "servers"
	resourcePath     = "tags"
)

func rootURL(c gophercloud.Client, serverID string) string {
	return c.ServiceURL(rootResourcePath, serverID, resourcePath)
}

func resourceURL(c gophercloud.Client, serverID, tag string) string {
	return c.ServiceURL(rootResourcePath, serverID, resourcePath, tag)
}

func listURL(c gophercloud.Client, serverID string) string {
	return rootURL(c, serverID)
}

func checkURL(c gophercloud.Client, serverID, tag string) string {
	return resourceURL(c, serverID, tag)
}

func replaceAllURL(c gophercloud.Client, serverID string) string {
	return rootURL(c, serverID)
}

func addURL(c gophercloud.Client, serverID, tag string) string {
	return resourceURL(c, serverID, tag)
}

func deleteURL(c gophercloud.Client, serverID, tag string) string {
	return resourceURL(c, serverID, tag)
}

func deleteAllURL(c gophercloud.Client, serverID string) string {
	return rootURL(c, serverID)
}
