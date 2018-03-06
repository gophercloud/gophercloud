package extradhcpopts

import "github.com/gophercloud/gophercloud"

// ExtraDHCPOptsExt is a struct that contains different DHCP options for a single port.
type ExtraDHCPOptsExt struct {
	ExtraDHCPOpts []ExtraDHCPOpts `json:"extra_dhcp_opts"`
}

// ExtraDHCPOpts represents a single set of extra DHCP options for a single port.
type ExtraDHCPOpts struct {
	// Name is the name of a single DHCP option.
	OptName string `json:"opt_name"`

	// Value is the value of a single DHCP option.
	OptValue string `json:"opt_value"`

	// IPVersion is the IP protocol version of a single DHCP option.
	// Valid value is 4 or 6. Default is 4.
	IPVersion int `json:"ip_version,omitempty"`
}

// ToMap is a helper function to convert an individual DHCPOpts structure
// into a sub-map.
func (opts ExtraDHCPOpts) ToMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}
