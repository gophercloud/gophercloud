package resourcetypes

import "github.com/gophercloud/gophercloud/v2"

const (
	resTypesPath = "resource_types"
)

func listURL(c gophercloud.Client) string {
	return c.ServiceURL(resTypesPath)
}

func getSchemaURL(c gophercloud.Client, resourceType string) string {
	return c.ServiceURL(resTypesPath, resourceType)
}

func generateTemplateURL(c gophercloud.Client, resourceType string) string {
	return c.ServiceURL(resTypesPath, resourceType, "template")
}
