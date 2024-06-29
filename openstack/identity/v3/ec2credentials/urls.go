package ec2credentials

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2")
}

func getURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2", id)
}

func createURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2")
}

func deleteURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2", id)
}
