package networks

type NetworkProvider struct {
	ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	ProviderNetworkType     string `json:"provider:network_type"`
}

type Network struct {
	Status       string        `json:"status"`
	Subnets      []interface{} `json:"subnets"`
	Name         string        `json:"name"`
	AdminStateUp bool          `json:"admin_state_up"`
	TenantID     string        `json:"tenant_id"`
	Shared       bool          `json:"shared"`
	ID           string        `json:"id"`
}

type NetworkResult struct {
	Network
	NetworkProvider
	RouterExternal bool `json:"router:external"`
}

type NetworkCreateResult struct {
	Network
	Segments            []NetworkProvider `json:"segments"`
	PortSecurityEnabled bool              `json:"port_security_enabled"`
}
