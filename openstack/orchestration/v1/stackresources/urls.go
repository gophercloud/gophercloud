package stackresources

import "github.com/gophercloud/gophercloud/v2"

func findURL(c gophercloud.Client, stackName string) string {
	return c.ServiceURL("stacks", stackName, "resources")
}

func listURL(c gophercloud.Client, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources")
}

func getURL(c gophercloud.Client, stackName, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources", resourceName)
}

func metadataURL(c gophercloud.Client, stackName, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources", resourceName, "metadata")
}

func listTypesURL(c gophercloud.Client) string {
	return c.ServiceURL("resource_types")
}

func schemaURL(c gophercloud.Client, typeName string) string {
	return c.ServiceURL("resource_types", typeName)
}

func templateURL(c gophercloud.Client, typeName string) string {
	return c.ServiceURL("resource_types", typeName, "template")
}

func markUnhealthyURL(c gophercloud.Client, stackName, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources", resourceName)
}
