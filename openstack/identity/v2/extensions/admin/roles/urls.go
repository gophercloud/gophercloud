package roles

import "github.com/rackspace/gophercloud"

const (
	extPath  = "OS-KSADMN"
	rolePath = "roles"
)

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(extPath, rolePath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(extPath, rolePath)
}
