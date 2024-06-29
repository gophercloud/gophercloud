package allocations

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("allocations")
}

func listURL(client gophercloud.Client) string {
	return createURL(client)
}

func resourceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("allocations", id)
}

func deleteURL(client gophercloud.Client, id string) string {
	return resourceURL(client, id)
}

func getURL(client gophercloud.Client, id string) string {
	return resourceURL(client, id)
}
