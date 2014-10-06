package provider

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
)

type AdminState *bool

// Convenience vars for AdminStateUp values.
var (
	iTrue  = true
	iFalse = false

	Nothing AdminState = nil
	Up      AdminState = &iTrue
	Down    AdminState = &iFalse
)

// NetworkExtAttrs represents an extended form of a Network with additional fields.
type NetworkExtAttrs struct {
	// UUID for the network
	ID string `mapstructure:"id" json:"id"`

	// Human-readable name for the network. Might not be unique.
	Name string `mapstructure:"name" json:"name"`

	// The administrative state of network. If false (down), the network does not forward packets.
	AdminStateUp bool `mapstructure:"admin_state_up" json:"admin_state_up"`

	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.
	Status string `mapstructure:"status" json:"status"`

	// Subnets associated with this network.
	Subnets []string `mapstructure:"subnets" json:"subnets"`

	// Owner of network. Only admin users can specify a tenant_id other than its own.
	TenantID string `mapstructure:"tenant_id" json:"tenant_id"`

	// Specifies whether the network resource can be accessed by any tenant or not.
	Shared bool `mapstructure:"shared" json:"shared"`

	// Specifies the nature of the physical network mapped to this network
	// resource. Examples are flat, vlan, or gre.
	NetworkType string `json:"provider:network_type" mapstructure:"provider:network_type"`

	// Identifies the physical network on top of which this network object is
	// being implemented. The OpenStack Networking API does not expose any facility
	// for retrieving the list of available physical networks. As an example, in
	// the Open vSwitch plug-in this is a symbolic name which is then mapped to
	// specific bridges on each compute host through the Open vSwitch plug-in
	// configuration file.
	PhysicalNetwork string `json:"provider:physical_network" mapstructure:"provider:physical_network"`

	// Identifies an isolated segment on the physical network; the nature of the
	// segment depends on the segmentation model defined by network_type. For
	// instance, if network_type is vlan, then this is a vlan identifier;
	// otherwise, if network_type is gre, then this will be a gre key.
	SegmentationID string `json:"provider:segmentation_id" mapstructure:"provider:segmentation_id"`
}

// ExtractGet decorates a GetResult struct returned from a networks.Get()
// function with extended attributes.
func ExtractGet(r networks.GetResult) (*NetworkExtAttrs, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	var res struct {
		Network *NetworkExtAttrs `json:"network"`
	}
	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron network: %v", err)
	}
	return res.Network, nil
}

// ExtractGet decorates a CreateResult struct returned from a networks.Create()
// function with extended attributes.
func ExtractCreate(r networks.CreateResult) (*NetworkExtAttrs, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	var res struct {
		Network *NetworkExtAttrs `json:"network"`
	}
	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron network: %v", err)
	}
	return res.Network, nil
}

// ExtractUpdate decorates a UpdateResult struct returned from a
// networks.Update() function with extended attributes.
func ExtractUpdate(r networks.UpdateResult) (*NetworkExtAttrs, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	var res struct {
		Network *NetworkExtAttrs `json:"network"`
	}
	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron network: %v", err)
	}
	return res.Network, nil
}

// ExtractList accepts a Page struct, specifically a NetworkPage struct, and
// extracts the elements into a slice of NetworkExtAttrs structs. In other
// words, a generic collection is mapped into a relevant slice.
func ExtractList(page pagination.Page) ([]NetworkExtAttrs, error) {
	var resp struct {
		Networks []NetworkExtAttrs `mapstructure:"networks" json:"networks"`
	}

	err := mapstructure.Decode(page.(networks.NetworkPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Networks, nil
}
