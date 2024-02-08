package resourceproviders

import "github.com/gophercloud/gophercloud/v2"

const (
	apiName = "resource_providers"
)

func resourceProvidersListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}

func deleteURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID)
}

func getURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID)
}

func updateURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID)
}

func getResourceProviderUsagesURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "usages")
}

func getResourceProviderInventoriesURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "inventories")
}

func getResourceProviderAllocationsURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "allocations")
}

func getResourceProviderTraitsURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "traits")
}
