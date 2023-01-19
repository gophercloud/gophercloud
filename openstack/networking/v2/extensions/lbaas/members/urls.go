package members

import "github.com/bizflycloud/gophercloud"

const (
	rootPath     = "lb"
	resourcePath = "members"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
