package webhooks

import "github.com/rackspace/gophercloud"

func listURL(c *gophercloud.ServiceClient, groupID, policyID string) string {
	return c.ServiceURL("groups", groupID, "policies", policyID, "webhooks")
}

func createURL(c *gophercloud.ServiceClient, groupID, policyID string) string {
	return c.ServiceURL("groups", groupID, "policies", policyID, "webhooks")
}

func getURL(c *gophercloud.ServiceClient, groupID, policyID, webhookID string) string {
	return c.ServiceURL("groups", groupID, "policies", policyID, "webhooks", webhookID)
}
