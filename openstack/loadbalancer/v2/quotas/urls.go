package quotas

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "quotas"

func resourceURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID)
}

func getURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceURL(c, projectID)
}

func updateURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceURL(c, projectID)
}
