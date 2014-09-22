package volumes

import "github.com/rackspace/gophercloud"

func volumesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("volumes")
}

func volumeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id)
}
