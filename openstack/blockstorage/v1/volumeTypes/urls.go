package volumeTypes

import "github.com/rackspace/gophercloud"

func volumeTypesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("types")
}

func volumeTypeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("types", id)
}
