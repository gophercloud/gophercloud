package tenants

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("tenants")
}

func getURL(client gophercloud.Client, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("tenants")
}

func deleteURL(client gophercloud.Client, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func updateURL(client gophercloud.Client, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}
