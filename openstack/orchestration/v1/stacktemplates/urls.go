package stacktemplates

import "github.com/gophercloud/gophercloud/v2"

func getURL(c gophercloud.Client, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}

func validateURL(c gophercloud.Client) string {
	return c.ServiceURL("validate")
}
