package speakers

import "github.com/gophercloud/gophercloud/v2"

const urlBase = "bgp-speakers"

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgp-speakers
func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers
func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgp-speakers
func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgp-speakers/{bgp-peer-id}
func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers/{bgp-peer-id}
func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/add_bgp_peer
func addBGPPeerURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "add_bgp_peer")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/remove_bgp_peer
func removeBGPPeerURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "remove_bgp_peer")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/get_advertised_routes
func getAdvertisedRoutesURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "get_advertised_routes")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/add_gateway_network
func addGatewayNetworkURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "add_gateway_network")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/remove_gateway_network
func removeGatewayNetworkURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "remove_gateway_network")
}
