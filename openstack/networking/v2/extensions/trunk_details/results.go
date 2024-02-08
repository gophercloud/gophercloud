package trunk_details

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
)

// TrunkDetailsExt represents additional trunking information returned in a
// ports query.
type TrunkDetailsExt struct {
	// trunk_details contains details of any trunk associated with the port
	TrunkDetails `json:"trunk_details,omitempty"`
}

// TrunkDetails contains additional trunking information returned in a
// ports query.
type TrunkDetails struct {
	// trunk_id contains the UUID of the trunk
	TrunkID string `json:"trunk_id"`

	// sub_ports contains a list of subports associated with the trunk
	SubPorts []Subport `json:"sub_ports,omitempty"`
}

type Subport struct {
	trunks.Subport

	// mac_address contains the MAC address of the subport.
	// Note that MACAddress may not be returned in list queries
	MACAddress string `json:"mac_address,omitempty"`
}
