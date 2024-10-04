package bgpvpns

import "github.com/gophercloud/gophercloud/v2"

const urlBase = "bgpvpn/bgpvpns"

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}
func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgpvpn/bgpvpns
func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}
func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgpvpn/bgpvpns
func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgpvpn/bgpvpns
func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}
func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}
func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations
func networkAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "network_associations")
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations/{network-association-id}
func networkAssociationResourceURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "network_associations", id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations
func listNetworkAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return networkAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations
func createNetworkAssociationURL(c gophercloud.Client, bgpVpnID string) string {
	return networkAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations/{network-association-id}
func getNetworkAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return networkAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/network_associations/{network-association-id}
func deleteNetworkAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return networkAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations
func routerAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "router_associations")
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations/{router-association-id}
func routerAssociationResourceURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "router_associations", id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations
func listRouterAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return routerAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations
func createRouterAssociationURL(c gophercloud.Client, bgpVpnID string) string {
	return routerAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations/{router-association-id}
func getRouterAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return routerAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations/{router-association-id}
func updateRouterAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return routerAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/router_associations/{router-association-id}
func deleteRouterAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return routerAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations
func portAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "port_associations")
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations/{port-association-id}
func portAssociationResourceURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return c.ServiceURL(urlBase, bgpVpnID, "port_associations", id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations
func listPortAssociationsURL(c gophercloud.Client, bgpVpnID string) string {
	return portAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations
func createPortAssociationURL(c gophercloud.Client, bgpVpnID string) string {
	return portAssociationsURL(c, bgpVpnID)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations/{port-association-id}
func getPortAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return portAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations/{port-association-id}
func updatePortAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return portAssociationResourceURL(c, bgpVpnID, id)
}

// return /v2.0/bgpvpn/bgpvpns/{bgpvpn-id}/port_associations/{port-association-id}
func deletePortAssociationURL(c gophercloud.Client, bgpVpnID string, id string) string {
	return portAssociationResourceURL(c, bgpVpnID, id)
}
