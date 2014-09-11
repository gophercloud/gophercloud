package networks

// A Network represents a a virtual layer-2 broadcast domain.
type Network struct {
	// Id is the unique identifier for the network.
	Id string `json:"id"`
	// Name is the (not necessarily unique) human-readable identifier for the network.
	Name string `json:"name"`
	// AdminStateUp is administrative state of the network. If false, network is down.
	AdminStateUp bool `json:"admin_state_up"`
	// Status indicates if the network is operational. Possible values: active, down, build, error.
	Status string `json:"status"`
	// Subnets are IP address blocks that can be used to assign IP addresses to virtual instances.
	Subnets []string `json:"subnets"`
	// Shared indicates whether the network can be accessed by any tenant or not.
	Shared bool `json:"shared"`
	// TenantId is the owner of the network. Admins may specify TenantId other than their own.
	TenantId string `json:"tenant_id"`
	// RouterExternal indicates if the network is connected to an external router.
	RouterExternal bool `json:"router:external"`
	// ProviderPhysicalNetwork is the name of the provider physical network.
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	// ProviderNetworkType is the type of provider network (eg "vlan").
	ProviderNetworkType string `json:"provider:network_type"`
	// ProviderSegmentationId is the provider network identifier (such as the vlan id).
	ProviderSegmentationId string `json:"provider:segmentation_id"`
}
