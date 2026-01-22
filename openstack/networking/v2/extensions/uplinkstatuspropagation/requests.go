package uplinkstatuspropagation

import (
	"net/url"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

// PortPropagateUplinkStatusListOptsExt adds uplink status propagation options
// to the base ports.ListOpts.
type PortPropagateUplinkStatusListOptsExt struct {
	ports.ListOptsBuilder

	PropagateUplinkStatus *bool `q:"propagate_uplink_status"`
}

// ToPortListQuery adds the propagate_uplink_status option to the base port
// list options.
func (opts PortPropagateUplinkStatusListOptsExt) ToPortListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts.ListOptsBuilder)
	if err != nil {
		return "", err
	}

	params := q.Query()

	if opts.PropagateUplinkStatus != nil {
		v := strconv.FormatBool(*opts.PropagateUplinkStatus)
		params.Add("propagate_uplink_status", v)
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

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
