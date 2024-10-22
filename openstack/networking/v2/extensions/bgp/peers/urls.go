package peers

import "github.com/gophercloud/gophercloud/v2"

const urlBase = "bgp-peers"

// return /v2.0/bgp-peers/{bgp-peer-id}
func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgp-peers
func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-peers
func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgp-peers
func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-peers/{bgp-peer-id}
func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}
