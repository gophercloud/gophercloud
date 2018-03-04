package extradhcpopts

// DHCPOptsExt is a struct that contains different DHCP options for a single port.
type DHCPOptsExt struct {
	DHCPOpts []DHCPOpts `json:"extra_dhcp_opts"`
}

// DHCPOpts represents a single set of extra DHCP options for a single port.
type DHCPOpts struct {
	// Name is the name of a single DHCP option.
	DHCPOptName string `json:"opt_name"`

	// Value is the value of a single DHCP option.
	DHCPOptValue string `json:"opt_value"`

	// IPVersion is the IP protocol version of a single DHCP option.
	// Valid value is 4 or 6. Default is 4.
	DHCPOptIPVersion int `json:"ip_version,omitempty"`
}
