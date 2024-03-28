package usage

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-simple-tenant-usage"

func allTenantsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func getTenantURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL(resourcePath, tenantID)
}
