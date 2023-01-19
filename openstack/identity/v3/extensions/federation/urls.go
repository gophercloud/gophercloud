package federation

import "github.com/bizflycloud/gophercloud"

const (
	rootPath     = "OS-FEDERATION"
	mappingsPath = "mappings"
)

func mappingsRootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, mappingsPath)
}
