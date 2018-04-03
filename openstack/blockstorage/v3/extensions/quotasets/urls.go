package quotasets

import "github.com/gophercloud/gophercloud"

const resourcePath = "os-quota-sets"

func resourceURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func getURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID)
}

func getDetailURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID+"?usage=true")
}

func updateURL(c *gophercloud.ServiceClient, projectID string) string {
	return getURL(c, projectID)
}

func deleteURL(c *gophercloud.ServiceClient, projectID string) string {
	return getURL(c, projectID)
}
