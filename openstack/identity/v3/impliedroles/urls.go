package impliedroles

import "github.com/gophercloud/gophercloud"

const (
	rolePath = "roles"
)

func listURL(client *gophercloud.ServiceClient, priorRoleId string) string {
	return client.ServiceURL(rolePath, priorRoleId, "implies")
}

func createURL(client *gophercloud.ServiceClient, priorRoleId string, impliesRoleId string) string {
	return client.ServiceURL(rolePath, priorRoleId, "implies", impliesRoleId)
}

func deleteURL(client *gophercloud.ServiceClient, priorRoleId string, impliesRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleId, "implies", impliesRoleID)
}
