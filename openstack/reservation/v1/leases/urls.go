package leases

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("leases")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("leases", id)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("leases")
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("leases", id)
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("leases", id)
}
