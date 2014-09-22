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
	// UUID representing the subnet
	ID string `mapstructure:"id" json:"id"`
	// UUID of the parent network
	NetworkID string `mapstructure:"network_id" json:"network_id"`
	// Human-readable name for the subnet. Might not be unique.
	Name string `mapstructure:"name" json:"name"`
	// IP version, either `4' or `6'
	IPVersion int `mapstructure:"ip_version" json:"ip_version"`
	// CIDR representing IP range for this subnet, based on IP version
	CIDR string `mapstructure:"cidr" json:"cidr"`
	// Default gateway used by devices in this subnet
	GatewayIP string `mapstructure:"gateway_ip" json:"gateway_ip"`
	// DNS name servers used by hosts in this subnet.
	DNSNameservers []string `mapstructure:"dns_nameservers" json:"dns_nameservers"`
	// Sub-ranges of CIDR available for dynamic allocation to ports. See AllocationPool.
	AllocationPools []AllocationPool `mapstructure:"allocation_pools" json:"allocation_pools"`
	// Routes that should be used by devices with IPs from this subnet (not including local subnet route).
	HostRoutes []interface{} `mapstructure:"host_routes" json:"host_routes"`
	// Specifies whether DHCP is enabled for this subnet or not.
	EnableDHCP bool `mapstructure:"enable_dhcp" json:"enable_dhcp"`
	// Owner of network. Only admin users can specify a tenant_id other than its own.
	TenantID string `mapstructure:"tenant_id" json:"tenant_id"`
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
