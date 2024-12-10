package orders

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("orders")
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("orders", id)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("orders")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("orders", id)
}
