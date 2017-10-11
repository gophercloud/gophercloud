package simpletenantusage

import "code.comcast.com/onecloud/gophercloud"

const resourcePath = "os-simple-tenant-usage"

func getURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func getTenantURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL(resourcePath, tenantID)
}
