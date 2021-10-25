package impliedroles

import "github.com/gophercloud/gophercloud"

const (
	rolePath = "roles"
)

func getURL(client *gophercloud.ServiceClient, priorRoleId string) string {
	return client.ServiceURL(rolePath, priorRoleId, "implies")
}
