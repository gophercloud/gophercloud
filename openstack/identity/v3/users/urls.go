package users

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("users")
}

func getURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("users")
}

func updateURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID)
}

func changePasswordURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "password")
}

func deleteURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID)
}

func listGroupsURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "groups")
}

func addToGroupURL(client gophercloud.Client, groupID, userID string) string {
	return client.ServiceURL("groups", groupID, "users", userID)
}

func isMemberOfGroupURL(client gophercloud.Client, groupID, userID string) string {
	return client.ServiceURL("groups", groupID, "users", userID)
}

func removeFromGroupURL(client gophercloud.Client, groupID, userID string) string {
	return client.ServiceURL("groups", groupID, "users", userID)
}

func listProjectsURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "projects")
}

func listInGroupURL(client gophercloud.Client, groupID string) string {
	return client.ServiceURL("groups", groupID, "users")
}
