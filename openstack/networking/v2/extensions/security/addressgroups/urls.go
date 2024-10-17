package addressgroups

import "github.com/gophercloud/gophercloud/v2"

const rootPath = "address-groups"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func resourceAddAddressesURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "add_addresses")
}

func resourceRemoveAddressesURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "remove_addresses")
}
