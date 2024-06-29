package federation

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "OS-FEDERATION"
	mappingsPath = "mappings"
)

func mappingsRootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, mappingsPath)
}

func mappingsResourceURL(c gophercloud.Client, mappingID string) string {
	return c.ServiceURL(rootPath, mappingsPath, mappingID)
}
