package ports

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("ports")
}

func listURL(client gophercloud.Client) string {
	return createURL(client)
}

func listDetailURL(client gophercloud.Client) string {
	return client.ServiceURL("ports", "detail")
}

func resourceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("ports", id)
}

func deleteURL(client gophercloud.Client, id string) string {
	return resourceURL(client, id)
}

func getURL(client gophercloud.Client, id string) string {
	return resourceURL(client, id)
}

func updateURL(client gophercloud.Client, id string) string {
	return resourceURL(client, id)
}
