package networks

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type NetworkProvider struct {
	ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	ProviderNetworkType     string `json:"provider:network_type"`
}

type Network struct {
	Status                  string        `mapstructure:"status" json:"status"`
	Subnets                 []interface{} `mapstructure:"subnets" json:"subnets"`
	Name                    string        `mapstructure:"name" json:"name"`
	AdminStateUp            bool          `mapstructure:"admin_state_up" json:"admin_state_up"`
	TenantID                string        `mapstructure:"tenant_id" json:"tenant_id"`
	Shared                  bool          `mapstructure:"shared" json:"shared"`
	ID                      string        `mapstructure:"id" json:"id"`
	ProviderSegmentationID  int           `mapstructure:"provider:segmentation_id" json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string        `mapstructure:"provider:physical_network" json:"provider:physical_network"`
	ProviderNetworkType     string        `mapstructure:"provider:network_type" json:"provider:network_type"`
	RouterExternal          bool          `mapstructure:"router:external" json:"router:external"`
}

type NetworkCreateResult struct {
	Network
	Segments            []NetworkProvider `json:"segments"`
	PortSecurityEnabled bool              `json:"port_security_enabled"`
}

func ExtractNetworks(page gophercloud.Page) ([]Network, error) {
	var resp struct {
		Networks []Network `mapstructure:"networks" json:"networks"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Networks, nil
}
