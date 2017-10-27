package roles

import "github.com/gophercloud/gophercloud"

const (
	rolePath    = "roles"
	userPath    = "users"
	domainPath  = "domains"
	groupPath   = "groups"
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

func groupOnDomainURL(client *gophercloud.ServiceClient, domainID, groupID, roleID string) string {
	return client.ServiceURL(domainPath, domainID, groupPath, groupID, rolePath, roleID)
}

func userOnDomainURL(client *gophercloud.ServiceClient, domainID, userID, roleID string) string {
	return client.ServiceURL(domainPath, domainID, userPath, userID, rolePath, roleID)
}

func groupOnProjectURL(client *gophercloud.ServiceClient, projectID, groupID, roleID string) string {
	return client.ServiceURL(projectPath, projectID, groupPath, groupID, rolePath, roleID)
}

func userOnProjectURL(client *gophercloud.ServiceClient, projectID, userID, roleID string) string {
	return client.ServiceURL(projectPath, projectID, userPath, userID, rolePath, roleID)
}
