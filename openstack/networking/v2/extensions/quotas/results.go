package quotas

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type Quota struct {
	Subnet        int `json:"subnet"`
	Router        int `json:"router"`
	Port          int `json:"port"`
	Network       int `json:"network"`
	FloatingIP    int `json:"floatingip"`
	HealthMonitor int `json:"health_monitor"`
	SecGroupRule  int `json:"security_group_rule"`
	SecGroup      int `json:"security_group"`
	VIP           int `json:"vip"`
	Member        int `json:"member"`
	Pool          int `json:"pool"`
}

type commonResult gophercloud.CommonResult

func (r commonResult) Extract() (*Quota, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Quota *Quota `json:"quota"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron quotas: %v", err)
	}

	return res.Quota, nil
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult commonResult
