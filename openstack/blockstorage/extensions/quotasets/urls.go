package quotasets

import "github.com/gophercloud/gophercloud"

const resourcePath = "os-quota-sets"

func getURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID)
}
