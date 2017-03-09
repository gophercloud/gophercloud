package users

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("users")
}

func getURL(client *gophercloud.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}
