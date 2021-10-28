package speakers

import "github.com/gophercloud/gophercloud"

const urlBase = "bgp-speakers"

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgp-speakers
func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers
func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
