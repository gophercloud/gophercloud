package users

import "github.com/gophercloud/gophercloud"

func baseURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("users")
}

func userURL(client *gophercloud.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}
