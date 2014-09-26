package floatingips

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// FloatingIP represents a floating IP resource. A floating IP is an external
// IP address that is mapped to an internal port and, optionally, a specific
// IP address on a private network. In other words, it enables access to an
// instance on a private network from an external network. For thsi reason,
// floating IPs can only be defined on networks where the `router:external'
// attribute (provided by the external network extension) is set to True.
type FloatingIP struct {
	// Unique identifier for the floating IP instance.
	ID string `json:"id" mapstructure:"id"`

	// UUID of the external network where the floating IP is to be created.
	FloatingNetworkID string `json:"floating_network_id" mapstructure:"floating_network_id"`

	// Address of the floating IP on the external network.
	FloatingIP string `json:"floating_ip_address" mapstructure:"floating_ip_address"`

	// UUID of the port on an internal network that is associated with the floating IP.
	PortID string `json:"port_id" mapstructure:"port_id"`

	// The specific IP address of the internal port which should be associated
	// with the floating IP.
	FixedIP string `json:"fixed_ip_address" mapstructure:"fixed_ip_address"`

	// Owner of the floating IP. Only admin users can specify a tenant identifier
	// other than its own.
	TenantID string `json:"tenant_id" mapstructure:"tenant_id"`

	// The condition of the API resource.
	Status string `json:"status" mapstructure:"status"`
}

type commonResult struct {
	gophercloud.CommonResult
}

// Extract a result and extracts a FloatingIP resource.
func (r commonResult) Extract() (*FloatingIP, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		FloatingIP *FloatingIP `json:"floatingip"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron floating IP: %v", err)
	}

	return res.FloatingIP, nil
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of an update operation.
type DeleteResult commonResult

type FloatingIPPage struct {
	pagination.LinkedPageBase
}

func (p FloatingIPPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"floatingips_links"`
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

func (p FloatingIPPage) IsEmpty() (bool, error) {
	is, err := ExtractFloatingIPs(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractFloatingIPs(page pagination.Page) ([]FloatingIP, error) {
	var resp struct {
		FloatingIPs []FloatingIP `mapstructure:"floatingips" json:"floatingips"`
	}

	err := mapstructure.Decode(page.(FloatingIPPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.FloatingIPs, nil
}
