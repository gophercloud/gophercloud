package usage

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-simple-tenant-usage"

func allTenantsURL(client gophercloud.Client) string {
	return client.ServiceURL(resourcePath)
}

func getTenantURL(client gophercloud.Client, tenantID string) string {
	return client.ServiceURL(resourcePath, tenantID)
}
