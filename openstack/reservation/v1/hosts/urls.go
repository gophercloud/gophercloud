package hosts

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("os-hosts")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("os-hosts", id)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("os-hosts")
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("os-hosts", id)
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("os-hosts", id)
}

func listAllocationsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("os-hosts", "allocations")
}

func getAllocationURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("os-hosts", id, "allocation")
}

func listResourcePropertiesURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("os-hosts", "properties")
}

func updateResourcePropertyURL(client *gophercloud.ServiceClient, name string) string {
	return client.ServiceURL("os-hosts", "properties", name)
}
