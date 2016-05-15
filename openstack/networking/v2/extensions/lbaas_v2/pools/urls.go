package pools

import "github.com/rackspace/gophercloud"

const (
	rootPath     = "lbaas"
	resourcePath = "pools"
	memeberPath  = "members"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func memberRootURL(c *gophercloud.ServiceClient, poolId string) string {
	return c.ServiceURL(rootPath, resourcePath, poolId, memeberPath)
}

func memberResourceURL(c *gophercloud.ServiceClient, poolID string, memeberID string) string {
	return c.ServiceURL(rootPath, resourcePath, poolID, memeberPath, memeberID)
}
