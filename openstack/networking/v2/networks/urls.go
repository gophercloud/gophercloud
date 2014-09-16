package networks

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func NetworkURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "networks", id)
}

func CreateURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "networks")
}
