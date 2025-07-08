package leases

import "github.com/gophercloud/gophercloud/v2"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("leases")
}

func listURL(client *gophercloud.ServiceClient) string {
	return createURL(client)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("leases", id)
}
