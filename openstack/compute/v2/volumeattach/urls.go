package volumeattach

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-volume_attachments"

func resourceURL(c gophercloud.Client, serverID string) string {
	return c.ServiceURL("servers", serverID, resourcePath)
}

func listURL(c gophercloud.Client, serverID string) string {
	return resourceURL(c, serverID)
}

func createURL(c gophercloud.Client, serverID string) string {
	return resourceURL(c, serverID)
}

func getURL(c gophercloud.Client, serverID, aID string) string {
	return c.ServiceURL("servers", serverID, resourcePath, aID)
}

func deleteURL(c gophercloud.Client, serverID, aID string) string {
	return getURL(c, serverID, aID)
}
