package extradhcpopts

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

// CreateOptsExt adds port DHCP options to the base ports.CreateOpts.
type CreateOptsExt struct {
	// CreateOptsBuilder is the interface options structs have to satisfy in order
	// to be used in the main Create operation in this package.
	ports.CreateOptsBuilder

	// DHCPOpts field is a set of DHCP options for a single port.
	ExtraDHCPOpts []ExtraDHCPOpts `json:"extra_dhcp_opts,omitempty"`
}

// ToPortCreateMap casts a CreateOptsExt struct to a map.
func (opts CreateOptsExt) ToPortCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToPortCreateMap()
	if err != nil {
		return nil, err
	}

	port := base["port"].(map[string]interface{})

	// Convert opts.DHCPOpts to a slice of maps.
	if opts.ExtraDHCPOpts != nil {
		extraDHCPOpts := make([]map[string]interface{}, len(opts.ExtraDHCPOpts))
		for i, opt := range opts.ExtraDHCPOpts {
			extraDHCPOptMap, err := opt.ToMap()
			if err != nil {
				return nil, err
			}
			extraDHCPOpts[i] = extraDHCPOptMap
		}
		port["extra_dhcp_opts"] = extraDHCPOpts
	}

	return base, nil
}
