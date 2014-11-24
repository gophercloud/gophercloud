package defsecrules

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

const rulepath = "os-security-group-default-rules"

func resourceURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL(rulepath, strconv.Itoa(id))
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rulepath)
}
