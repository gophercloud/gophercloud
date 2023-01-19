package stacktemplates

import "github.com/bizflycloud/gophercloud"

func getURL(c *gophercloud.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}

func validateURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("validate")
}
