package subnets

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func ResourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "subnets", id)
}

func RootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "subnets")
}

func ListURL(c *gophercloud.ServiceClient) string {
	return RootURL(c)
}

func GetURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}
