package projects

import "github.com/gophercloud/gophercloud/v2"

func listAvailableURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("auth", "projects")
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("projects")
}

func getURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("projects")
}

func deleteURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func listTagsURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}

func modifyTagsURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}

func deleteTagsURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}

func getTagURL(client *gophercloud.ServiceClient, projectID string, tag string) string {
	return client.ServiceURL("projects", projectID, "tags", tag)
}

func addTagURL(client *gophercloud.ServiceClient, projectID string, tag string) string {
	return client.ServiceURL("projects", projectID, "tags", tag)
}

func deleteTagURL(client *gophercloud.ServiceClient, projectID string, tag string) string {
	return client.ServiceURL("projects", projectID, "tags", tag)
}
