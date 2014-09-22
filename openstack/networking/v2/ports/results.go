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
	// UUID for the port.
	ID string `mapstructure:"id" json:"id"`
	// Network that this port is associated with.
	NetworkID string `mapstructure:"network_id" json:"network_id"`
	// Human-readable name for the port. Might not be unique.
	Name string `mapstructure:"name" json:"name"`
	// Administrative state of port. If false (down), port does not forward packets.
	AdminStateUp bool `mapstructure:"admin_state_up" json:"admin_state_up"`
	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.
	Status string `mapstructure:"status" json:"status"`
	// Mac address to use on this port.
	MACAddress string `mapstructure:"mac_address" json:"mac_address"`
	// Specifies IP addresses for the port thus associating the port itself with
	// the subnets where the IP addresses are picked from
	FixedIPs []IP `mapstructure:"fixed_ips" json:"fixed_ips"`
	// Owner of network. Only admin users can specify a tenant_id other than its own.
	TenantID string `mapstructure:"tenant_id" json:"tenant_id"`
	// Identifies the entity (e.g.: dhcp agent) using this port.
	DeviceOwner string `mapstructure:"device_owner" json:"device_owner"`
	// Specifies the IDs of any security groups associated with a port.
	SecurityGroups []string `mapstructure:"security_groups" json:"security_groups"`
	// Identifies the device (e.g., virtual server) using this port.
	DeviceID string `mapstructure:"device_id" json:"device_id"`

	AllowedAddressPairs []interface{} `mapstructure:"allowed" json:"allowed"`
	ExtraDHCPOpts       interface{}   `mapstructure:"extra_dhcp_opts" json:"extra_dhcp_opts"`
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
