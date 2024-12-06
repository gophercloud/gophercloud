package request

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "zones"
	tasksPath    = "tasks"
	resourcePath = "transfer_requests"
)

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath)
}

func createURL(c gophercloud.Client, zoneID string) string {
	return c.ServiceURL(rootPath, zoneID, tasksPath, resourcePath)
}

func resourceURL(c gophercloud.Client, transferID string) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath, transferID)
}
