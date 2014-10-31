package lb

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

const path = "loadbalancers"

func resourceURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL(path, strconv.Itoa(id))
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(path)
}
