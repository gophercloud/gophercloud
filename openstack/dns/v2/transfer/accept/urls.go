package accept

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "zones"
	tasksPath    = "tasks"
	resourcePath = "transfer_accepts"
)

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath)
}

func resourceURL(c gophercloud.Client, transferAcceptID string) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath, transferAcceptID)
}
