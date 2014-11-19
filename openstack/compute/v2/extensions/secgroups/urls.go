package secgroups

import "github.com/rackspace/gophercloud"

const (
	secgrouppath = "os-security-groups"
)

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(secgrouppath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(secgrouppath)
}
