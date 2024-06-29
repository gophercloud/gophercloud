package applicationcredentials

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "application_credentials")
}

func getURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "application_credentials", id)
}

func createURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "application_credentials")
}

func deleteURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "application_credentials", id)
}

func listAccessRulesURL(client gophercloud.Client, userID string) string {
	return client.ServiceURL("users", userID, "access_rules")
}

func getAccessRuleURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "access_rules", id)
}

func deleteAccessRuleURL(client gophercloud.Client, userID string, id string) string {
	return client.ServiceURL("users", userID, "access_rules", id)
}
