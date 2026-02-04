package uplinkstatuspropagation

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

// PortPropagateUplinkStatusCreateOptsExt adds uplink status propagation
// options to the base ports.CreateOpts.
type PortPropagateUplinkStatusCreateOptsExt struct {
	ports.CreateOptsBuilder

	// PropagateUplinkStatus toggles propagation of uplink status to the
	// port.
	PropagateUplinkStatus *bool `json:"propagate_uplink_status,omitempty"`
}

// ToPortCreateMap casts a CreateOpts struct to a map.
func (opts PortPropagateUplinkStatusCreateOptsExt) ToPortCreateMap() (map[string]any, error) {
	base, err := opts.CreateOptsBuilder.ToPortCreateMap()
	if err != nil {
		return nil, err
	}

	port := base["port"].(map[string]any)

	if opts.PropagateUplinkStatus != nil {
		port["propagate_uplink_status"] = *opts.PropagateUplinkStatus
	}

	return base, nil
}

// PortPropagateUplinkStatusUpdateOptsExt adds uplink status propagation
// options to the base ports.UpdateOpts.
type PortPropagateUplinkStatusUpdateOptsExt struct {
	ports.UpdateOptsBuilder

	// PropagateUplinkStatus toggles propagation of uplink status to the
	// port.
	PropagateUplinkStatus *bool `json:"propagate_uplink_status,omitempty"`
}

// ToPortUpdateMap casts a UpdateOpts struct to a map.
func (opts PortPropagateUplinkStatusUpdateOptsExt) ToPortUpdateMap() (map[string]any, error) {
	base, err := opts.UpdateOptsBuilder.ToPortUpdateMap()
	if err != nil {
		return nil, err
	}

	port := base["port"].(map[string]any)

	if opts.PropagateUplinkStatus != nil {
		port["propagate_uplink_status"] = *opts.PropagateUplinkStatus
	}

	return base, nil
}
