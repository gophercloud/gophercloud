package providers

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "lbaas"
	resourcePath = "providers"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, resourcePath)
}
