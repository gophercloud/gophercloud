package sharenetworks

import "github.com/gophercloud/gophercloud"

// ShareNetwork contains all the information associated with an OpenStack
// ShareNetwork.
type ShareNetwork struct {
	// The Share Network ID
	ID string `json:"id"`
	// The UUID of the project where the share network was created
	ProjectID string `json:"project_id"`
	// The neutron network ID
	NeutronNetID string `json:"neutron_net_id"`
	// The neutron subnet ID
	NeutronSubnetID string `json:"neutron_subnet_id"`
	// The nova network ID
	NovaNetID string `json:"nova_net_id"`
	// The network type. A valid value is VLAN, VXLAN, GRE or flat
	NetworkType string `json:"network_type"`
	// The segmentation ID
	SegmentationID int `json:"segmentation_id"`
	// The IP block from which to allocate the network, in CIDR notation
	CIDR string `json:"cidr"`
	// The IP version of the network. A valid value is 4 or 6
	IPVersion int `json:"ip_version"`
	// The Share Network name
	Name string `json:"name"`
	// The Share Network description
	Description string `json:"description"`
	// The date and time stamp when the Share Network was created
	CreatedAt string `json:"created_at"`
	// The date and time stamp when the Share Network was updated
	UpdatedAt string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the ShareNetwork object out of the commonResult object.
func (r commonResult) Extract() (*ShareNetwork, error) {
	var s struct {
		ShareNetwork *ShareNetwork `json:"share_network"`
	}
	err := r.ExtractInto(&s)
	return s.ShareNetwork, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}
