package networks

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// NetworkProvider represents provider extension data
type NetworkProvider struct {
	ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	ProviderNetworkType     string `json:"provider:network_type"`
}

// Network represents, well, a network.
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

// NetworkCreateResult represents what is returned by a create operation.
type NetworkCreateResult struct {
	Network
	Segments            []NetworkProvider `json:"segments"`
	PortSecurityEnabled bool              `json:"port_security_enabled"`
}

// NetworkPage is the page returned by a pager when traversing over a
// collection of networks.
type NetworkPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of networks has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (p NetworkPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"networks_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
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

// IsEmpty checks whether a NetworkPage struct is empty.
func (p NetworkPage) IsEmpty() (bool, error) {
	is, err := ExtractNetworks(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractNetworks accepts a Page struct, specifically a NetworkPage struct,
// and extracts the elements into a slice of Network structs. In other words,
// a generic collection is mapped into a relevant slice.
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
