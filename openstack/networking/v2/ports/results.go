package ports

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type IP struct {
	SubnetID  string `mapstructure:"subnet_id" json:"subnet_id"`
	IPAddress string `mapstructure:"ip_address" json:"ip_address,omitempty"`
}

type Port struct {
	Status              string        `mapstructure:"status" json:"status"`
	Name                string        `mapstructure:"name" json:"name"`
	AllowedAddressPairs []interface{} `mapstructure:"allowed" json:"allowed"`
	AdminStateUp        bool          `mapstructure:"admin_state_up" json:"admin_state_up"`
	NetworkID           string        `mapstructure:"network_id" json:"network_id"`
	TenantID            string        `mapstructure:"tenant_id" json:"tenant_id"`
	ExtraDHCPOpts       interface{}   `mapstructure:"extra_dhcp_opts" json:"extra_dhcp_opts"`
	DeviceOwner         string        `mapstructure:"device_owner" json:"device_owner"`
	MACAddress          string        `mapstructure:"mac_address" json:"mac_address"`
	FixedIPs            []IP          `mapstructure:"fixed_ips" json:"fixed_ips"`
	ID                  string        `mapstructure:"id" json:"id"`
	SecurityGroups      []string      `mapstructure:"security_groups" json:"security_groups"`
	DeviceID            string        `mapstructure:"device_id" json:"device_id"`
	BindingHostID       string        `mapstructure:"binding:host_id" json:"binding:host_id"`
	BindingVIFDetails   interface{}   `mapstructure:"binding:vif_details" json:"binding:vif_details"`
	BindingVIFType      string        `mapstructure:"binding:vif_type" json:"binding:vif_type"`
	BindingProfile      interface{}   `mapstructure:"binding:profile" json:"binding:profile"`
	BindingVNICType     string        `mapstructure:"binding:vnic_type" json:"binding:vnic_type"`
}

type PortPage struct {
	pagination.LinkedPageBase
}

func (current PortPage) NextPageURL() (string, error) {
	type resp struct {
		Links []struct {
			Href string `mapstructure:"href"`
			Rel  string `mapstructure:"rel"`
		} `mapstructure:"ports_links"`
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

func (r PortPage) IsEmpty() (bool, error) {
	is, err := ExtractPorts(r)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractPorts(page pagination.Page) ([]Port, error) {
	var resp struct {
		Ports []Port `mapstructure:"ports" json:"ports"`
	}

	err := mapstructure.Decode(page.(PortPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Ports, nil
}
