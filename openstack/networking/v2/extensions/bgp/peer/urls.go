package peer

import "github.com/gophercloud/gophercloud"

// https://docs.openstack.org/api-ref/network/v2/#list-bgp-peers
const baseurl = "bgp-peers"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseurl, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseurl)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
