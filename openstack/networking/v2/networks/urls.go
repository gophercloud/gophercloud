package networks

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func ResourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "networks", id)
}

func RootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "networks")
}

func GetURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}

func ListURL(c *gophercloud.ServiceClient) string {
	return RootURL(c)
}

func CreateURL(c *gophercloud.ServiceClient) string {
	return RootURL(c)
}

func DeleteURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}
