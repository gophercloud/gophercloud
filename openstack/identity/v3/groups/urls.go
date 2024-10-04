package groups

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("groups")
}

func getURL(client gophercloud.Client, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("groups")
}

func updateURL(client gophercloud.Client, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func deleteURL(client gophercloud.Client, groupID string) string {
	return client.ServiceURL("groups", groupID)
}
