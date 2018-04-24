package usage

import "github.com/gophercloud/gophercloud"

const resourcePath = "os-simple-tenant-usage"

func getTenantURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL(resourcePath, tenantID)
}
