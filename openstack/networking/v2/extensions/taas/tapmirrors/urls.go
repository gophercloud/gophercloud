package tapmirrors

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "taas"
	resourcePath = "tap_mirrors"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}
