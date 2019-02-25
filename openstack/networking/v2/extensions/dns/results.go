package dns

// PortsBindingExt represents a decorated form of a Port with the additional
// port binding information.
type PortDNSExt struct {
	// The ID of the host where the port is allocated.
	DNSName string `json:"dns_name"`

	// The VIF type for the port.
	DNSAssignment []map[string]string `json:"dns_assignment"`
}

type FloatingIPDNSExt struct {
	// The ID of the host where the port is allocated.
	DNSName string `json:"dns_name"`

	// The VIF type for the port.
	DNSDomain string `json:"dns_domain"`
}

type NetworkDNSExt struct {
	// The VIF type for the port.
	DNSDomain string `json:"dns_domain"`
}
