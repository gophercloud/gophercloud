package agents

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "agents"
const dhcpNetworksResourcePath = "dhcp-networks"
const l3RoutersResourcePath = "l3-routers"
const bgpSpeakersResourcePath = "bgp-drinstances"
const bgpDRAgentSpeakersResourcePath = "bgp-speakers"
const bgpDRAgentAgentResourcePath = "bgp-dragents"

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func dhcpNetworksURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, dhcpNetworksResourcePath)
}

func l3RoutersURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, l3RoutersResourcePath)
}

func listDHCPNetworksURL(c gophercloud.Client, id string) string {
	return dhcpNetworksURL(c, id)
}

func listL3RoutersURL(c gophercloud.Client, id string) string {
	return l3RoutersURL(c, id)
}

func scheduleDHCPNetworkURL(c gophercloud.Client, id string) string {
	return dhcpNetworksURL(c, id)
}

func scheduleL3RouterURL(c gophercloud.Client, id string) string {
	return l3RoutersURL(c, id)
}

func removeDHCPNetworkURL(c gophercloud.Client, id string, networkID string) string {
	return c.ServiceURL(resourcePath, id, dhcpNetworksResourcePath, networkID)
}

func removeL3RouterURL(c gophercloud.Client, id string, routerID string) string {
	return c.ServiceURL(resourcePath, id, l3RoutersResourcePath, routerID)
}

// return /v2.0/agents/{agent-id}/bgp-drinstances
func listBGPSpeakersURL(c gophercloud.Client, agentID string) string {
	return c.ServiceURL(resourcePath, agentID, bgpSpeakersResourcePath)
}

// return /v2.0/agents/{agent-id}/bgp-drinstances
func scheduleBGPSpeakersURL(c gophercloud.Client, id string) string {
	return listBGPSpeakersURL(c, id)
}

// return /v2.0/agents/{agent-id}/bgp-drinstances/{bgp-speaker-id}
func removeBGPSpeakersURL(c gophercloud.Client, agentID string, speakerID string) string {
	return c.ServiceURL(resourcePath, agentID, bgpSpeakersResourcePath, speakerID)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/bgp-dragents
func listDRAgentHostingBGPSpeakersURL(c gophercloud.Client, speakerID string) string {
	return c.ServiceURL(bgpDRAgentSpeakersResourcePath, speakerID, bgpDRAgentAgentResourcePath)
}
