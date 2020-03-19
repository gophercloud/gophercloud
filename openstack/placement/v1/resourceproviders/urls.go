package resourceproviders

import "github.com/gophercloud/gophercloud"

const (
	apiName = "resource_providers"
)

func resourceProvidersListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}

func getResourceProviderUsagesURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "usages")
}

func getResourceProviderInventoriesURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "inventories")
}

func getResourceProviderTraitsURL(client *gophercloud.ServiceClient, resourceProviderID string) string {
	return client.ServiceURL(apiName, resourceProviderID, "traits")
}
