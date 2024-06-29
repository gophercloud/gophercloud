package quotasets

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-quota-sets"

func getURL(c gophercloud.Client, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID)
}

func getDetailURL(c gophercloud.Client, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID, "detail")
}

func updateURL(c gophercloud.Client, tenantID string) string {
	return getURL(c, tenantID)
}

func deleteURL(c gophercloud.Client, tenantID string) string {
	return getURL(c, tenantID)
}
