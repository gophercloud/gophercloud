package portsbinding

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a port resource.
func (r commonResult) Extract() (*Port, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Port *Port `json:"port"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Port, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// IP is a sub-struct that represents an individual IP.
type IP struct {
	SubnetID  string `mapstructure:"subnet_id" json:"subnet_id"`
	IPAddress string `mapstructure:"ip_address" json:"ip_address,omitempty"`
}

// Port represents a Neutron port. See package documentation for a top-level
// description of what this is.
type Port struct {
	// UUID for the port.
	ID string `mapstructure:"id" json:"id"`
	// Network that this port is associated with.
	NetworkID string `mapstructure:"network_id" json:"network_id"`
	// Human-readable name for the port. Might not be unique.
	Name string `mapstructure:"name" json:"name"`
	// Administrative state of port. If false (down), port does not forward packets.
	AdminStateUp bool `mapstructure:"admin_state_up" json:"admin_state_up"`
	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.
	Status string `mapstructure:"status" json:"status"`
	// Mac address to use on this port.
	MACAddress string `mapstructure:"mac_address" json:"mac_address"`
	// Specifies IP addresses for the port thus associating the port itself with
	// the subnets where the IP addresses are picked from
	FixedIPs []IP `mapstructure:"fixed_ips" json:"fixed_ips"`
	// Owner of network. Only admin users can specify a tenant_id other than its own.
	TenantID string `mapstructure:"tenant_id" json:"tenant_id"`
	// Identifies the entity (e.g.: dhcp agent) using this port.
	DeviceOwner string `mapstructure:"device_owner" json:"device_owner"`
	// Specifies the IDs of any security groups associated with a port.
	SecurityGroups []string `mapstructure:"security_groups" json:"security_groups"`
	// Identifies the device (e.g., virtual server) using this port.
	DeviceID string `mapstructure:"device_id" json:"device_id"`
	// The ID of the host where the port is allocated
	HostID string `mapstructure:"binding:host_id" json:"binding:host_id"`
	// The virtual network interface card (vNIC) type that is bound to the
	// neutron port
	VNICType string `mapstructure:"binding:vnic_type" json:"binding:vnic_type"`
	// A dictionary that enables the application running on the specified
	// host to pass and receive virtual network interface (VIF) port-specific
	// information to the plug-in
	Profile map[string]string `mapstructure:"binding:profile" json:"binding:profile"`
}
