package pools

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "lbaas"
	resourcePath = "pools"
	memberPath   = "members"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func memberRootURL(c gophercloud.Client, poolId string) string {
	return c.ServiceURL(rootPath, resourcePath, poolId, memberPath)
}

func memberResourceURL(c gophercloud.Client, poolID string, memberID string) string {
	return c.ServiceURL(rootPath, resourcePath, poolID, memberPath, memberID)
}
