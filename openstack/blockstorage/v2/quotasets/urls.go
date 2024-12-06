package quotasets

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-quota-sets"

func getURL(c gophercloud.Client, projectID string) string {
	return c.ServiceURL(resourcePath, projectID)
}

func getDefaultsURL(c gophercloud.Client, projectID string) string {
	return c.ServiceURL(resourcePath, projectID, "defaults")
}

func updateURL(c gophercloud.Client, projectID string) string {
	return getURL(c, projectID)
}

func deleteURL(c gophercloud.Client, projectID string) string {
	return getURL(c, projectID)
}
