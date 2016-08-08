package bays

import "github.com/gophercloud/gophercloud"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("bays")
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("bays")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("bays", id)
}
