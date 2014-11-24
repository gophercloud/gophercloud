package secgroups

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

const (
	secgrouppath = "os-security-groups"
	rulepath     = "os-security-group-rules"
)

func resourceURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL(secgrouppath, strconv.Itoa(id))
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(secgrouppath)
}

func listByServerURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL(secgrouppath, "servers", serverID, secgrouppath)
}

func rootRuleURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rulepath)
}

func resourceRuleURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL(rulepath, strconv.Itoa(id))
}

func serverActionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("servers", id, "action")
}
