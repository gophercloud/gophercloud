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

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("os-hosts", id)
}
