package routerinsertion

// FirewallGroupExt is an extension to the base Firewall group object
type FirewallGroupExt struct {
	// PortIDs are the ports that the firewall group is attached to.
	PortIDs []string `json:"ports"`
}
