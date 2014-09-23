package quotas

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Quota struct {
	Subnet     int `json:"subnet"`
	Router     int `json:"router"`
	Port       int `json:"port"`
	Network    int `json:"network"`
	FloatingIP int `json:"floatingip"`
}

type commonResult struct {
	resp map[string]interface{}
	err  error
}

type GetResult commonResult

func (r GetResult) Extract() (*Quota, error) {
	if r.err != nil {
		return nil, r.err
	}

	var res struct {
		Quota *Quota `json:"quota"`
	}

	err := mapstructure.Decode(r.resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron quotas: %v", err)
	}

	return res.Quota, nil
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*Quota, error) {
	if r.err != nil {
		return nil, r.err
	}

	var res struct {
		Quota *Quota `json:"quota"`
	}

	err := mapstructure.Decode(r.resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron quotas: %v", err)
	}

	return res.Quota, nil
}

type DeleteResult commonResult
