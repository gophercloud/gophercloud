package roles

import "github.com/gophercloud/gophercloud/v2"

const (
	rolePath = "roles"
)

func listURL(client gophercloud.Client) string {
	return client.ServiceURL(rolePath)
}

func getURL(client gophercloud.Client, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL(rolePath)
}

func updateURL(client gophercloud.Client, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func deleteURL(client gophercloud.Client, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func listAssignmentsURL(client gophercloud.Client) string {
	return client.ServiceURL("role_assignments")
}

func listAssignmentsOnResourceURL(client gophercloud.Client, targetType, targetID, actorType, actorID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath)
}

func assignURL(client gophercloud.Client, targetType, targetID, actorType, actorID, roleID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath, roleID)
}

func createRoleInferenceRuleURL(client gophercloud.Client, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}

func getRoleInferenceRuleURL(client gophercloud.Client, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}

func listRoleInferenceRulesURL(client gophercloud.Client) string {
	return client.ServiceURL("role_inferences")
}

func deleteRoleInferenceRuleURL(client gophercloud.Client, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}
