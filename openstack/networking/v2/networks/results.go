package networks

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type NetworkProvider struct {
	ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	ProviderNetworkType     string `json:"provider:network_type"`
}

type Network struct {
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

	ProviderSegmentationID  int    `mapstructure:"provider:segmentation_id" json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `mapstructure:"provider:physical_network" json:"provider:physical_network"`
	ProviderNetworkType     string `mapstructure:"provider:network_type" json:"provider:network_type"`
	RouterExternal          bool   `mapstructure:"router:external" json:"router:external"`
}

type NetworkCreateResult struct {
	Network
	Segments            []NetworkProvider `json:"segments"`
	PortSecurityEnabled bool              `json:"port_security_enabled"`
}

type NetworkPage struct {
	pagination.LinkedPageBase
}

func (current NetworkPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"networks_links"`
	}

	var r resp
	err := mapstructure.Decode(current.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

func (r NetworkPage) IsEmpty() (bool, error) {
	is, err := ExtractNetworks(r)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractNetworks(page pagination.Page) ([]Network, error) {
	var resp struct {
		Networks []Network `mapstructure:"networks" json:"networks"`
	}

	err := mapstructure.Decode(page.(NetworkPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Networks, nil
}
