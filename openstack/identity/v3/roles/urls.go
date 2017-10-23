package roles

import "github.com/gophercloud/gophercloud"

const (
	rolePath    = "roles"
	userPath    = "users"
	projectPath = "projects"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func getURL(client *gophercloud.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func listAssignmentsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("role_assignments")
}

func userOnProjectURL(client *gophercloud.ServiceClient, projectID, userID, roleID string) string {
	return client.ServiceURL(projectPath, projectID, userPath, userID, rolePath, roleID)
}
