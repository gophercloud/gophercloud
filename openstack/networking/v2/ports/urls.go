package ports

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func ResourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "ports", id)
}

func RootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "ports")
}

func ListURL(c *gophercloud.ServiceClient) string {
	return RootURL(c)
}

func GetURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}

func CreateURL(c *gophercloud.ServiceClient) string {
	return RootURL(c)
}

func UpdateURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}

func DeleteURL(c *gophercloud.ServiceClient, id string) string {
	return ResourceURL(c, id)
}
