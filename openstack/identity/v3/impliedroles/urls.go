package impliedroles

import "github.com/gophercloud/gophercloud"

const (
	rolePath = "roles"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("v3/role_inferences")
}

func createURL(client *gophercloud.ServiceClient, pirorRoleId string, impliesRoleId string) string {
	return client.ServiceURL(rolePath, pirorRoleId, "implies", impliesRoleId)
}

func deleteURL(client *gophercloud.ServiceClient, impliedRoleID string) string {
	return client.ServiceURL("v3/impliedrole", impliedRoleID)
}
