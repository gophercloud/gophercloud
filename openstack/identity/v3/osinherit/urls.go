package osinherit

import "github.com/gophercloud/gophercloud/v2"

const (
	inheritPath = "OS-INHERIT"
)

func assignURL(client *gophercloud.ServiceClient, targetType, targetID, actorType, actorID, roleID string) string {
	return client.ServiceURL(inheritPath, targetType, targetID, actorType, actorID, "roles", roleID, "inherited_to_projects")
}
