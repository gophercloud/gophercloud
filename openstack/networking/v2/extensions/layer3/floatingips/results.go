package floatingips

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type FloatingIP struct {
	ID                string `json:"id" mapstructure:"id"`
	FloatingNetworkID string `json:"floating_network_id" mapstructure:"floating_network_id"`
	FloatingIP        string `json:"floating_ip_address" mapstructure:"floating_ip_address"`
	PortID            string `json:"port_id" mapstructure:"port_id"`
	FixedIP           string `json:"fixed_ip_address" mapstructure:"fixed_ip_address"`
	TenantID          string `json:"tenant_id" mapstructure:"tenant_id"`
}

type commonResult struct {
	gophercloud.CommonResult
}

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

type CreateResult struct {
	commonResult
}
