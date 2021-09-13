package peer

import "github.com/gophercloud/gophercloud"

const urlBase = "bgp-peers"

// return /v2.0/bgp-peers/{bgp-peer-id}
func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgp-peers
func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-peers
func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

// return /v2.0/bgp-peers
func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
