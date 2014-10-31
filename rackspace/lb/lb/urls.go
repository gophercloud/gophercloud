package lb

import "github.com/rackspace/gophercloud"

const path = "loadbalancers"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(path, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(path)
}
