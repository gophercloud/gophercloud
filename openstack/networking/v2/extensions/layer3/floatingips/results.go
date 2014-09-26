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
	Status            string `json:"status" mapstructure:"status"`
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

type DeleteResult commonResult
