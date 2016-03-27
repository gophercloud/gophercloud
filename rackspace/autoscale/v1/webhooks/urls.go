package webhooks

import "github.com/rackspace/gophercloud"

func listURL(c *gophercloud.ServiceClient, groupID, policyID string) string {
	return c.ServiceURL("groups", groupID, "policies", policyID, "webhooks")
}
