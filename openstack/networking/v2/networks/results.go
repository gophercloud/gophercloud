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
