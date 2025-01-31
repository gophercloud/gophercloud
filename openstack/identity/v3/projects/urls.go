package projects

import "github.com/gophercloud/gophercloud/v2"

func listAvailableURL(client gophercloud.Client) string {
	return client.ServiceURL("auth", "projects")
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("projects")
}

func getURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("projects")
}

func deleteURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func listTagsURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}

func modifyTagsURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}

func deleteTagsURL(client gophercloud.Client, projectID string) string {
	return client.ServiceURL("projects", projectID, "tags")
}
