package subnets

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Subnet struct {
	Name            string           `mapstructure:"name" json:"name"`
	EnableDHCP      bool             `mapstructure:"enable_dhcp" json:"enable_dhcp"`
	NetworkID       string           `mapstructure:"network_id" json:"network_id"`
	TenantID        string           `mapstructure:"tenant_id" json:"tenant_id"`
	DNSNameservers  []interface{}    `mapstructure:"dns_nameservers" json:"dns_nameservers"`
	AllocationPools []AllocationPool `mapstructure:"allocation_pools" json:"allocation_pools"`
	HostRoutes      []interface{}    `mapstructure:"host_routes" json:"host_routes"`
	IPVersion       int              `mapstructure:"ip_version" json:"ip_version"`
	GatewayIP       string           `mapstructure:"gateway_ip" json:"gateway_ip"`
	CIDR            string           `mapstructure:"cidr" json:"cidr"`
	ID              string           `mapstructure:"id" json:"id"`
}

type SubnetPage struct {
	pagination.LinkedPageBase
}

func (current SubnetPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"subnets_links"`
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

func (r SubnetPage) IsEmpty() (bool, error) {
	is, err := ExtractSubnets(r)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractSubnets(page pagination.Page) ([]Subnet, error) {
	var resp struct {
		Subnets []Subnet `mapstructure:"subnets" json:"subnets"`
	}

	err := mapstructure.Decode(page.(SubnetPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Subnets, nil
}
