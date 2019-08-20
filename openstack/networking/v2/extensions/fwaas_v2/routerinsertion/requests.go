package routerinsertion

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
)

// CreateOptsExt adds a PortIDs option to the base CreateOpts.
type CreateOptsExt struct {
	firewall_groups.CreateOptsBuilder
	PortIDs []string `json:"ports"`
}

// ToFirewallGroupCreateMap adds ports to the base firewall creation options.
func (opts CreateOptsExt) ToFirewallGroupCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToFirewallGroupCreateMap()
	if err != nil {
		return nil, err
	}

	firewallMap := base["firewall_group"].(map[string]interface{})
	firewallMap["ports"] = opts.PortIDs

	return base, nil
}

// UpdateOptsExt updates a PortIDs option to the base UpdateOpts.
type UpdateOptsExt struct {
	firewall_groups.UpdateOptsBuilder
	PortIDs []string `json:"ports"`
}

// ToFirewallGroupUpdateMap adds ports to the base firewall update options.
func (opts UpdateOptsExt) ToFirewallGroupUpdateMap() (map[string]interface{}, error) {
	base, err := opts.UpdateOptsBuilder.ToFirewallGroupUpdateMap()
	if err != nil {
		return nil, err
	}

	firewallMap := base["firewall_group"].(map[string]interface{})
	firewallMap["ports"] = opts.PortIDs

	return base, nil
}
