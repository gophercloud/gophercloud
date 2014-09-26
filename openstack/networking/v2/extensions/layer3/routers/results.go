package routers

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type GatewayInfo struct {
	NetworkID string `json:"network_id" mapstructure:"network_id"`
}

type Router struct {
	Status       string      `json:"status" mapstructure:"status"`
	GatewayInfo  GatewayInfo `json:"external_gateway_info" mapstructure:"external_gateway_info"`
	AdminStateUp bool        `json:"admin_state_up" mapstructure:"admin_state_up"`
	Name         string      `json:"name" mapstructure:"name"`
	ID           string      `json:"id" mapstructure:"id"`
	TenantID     string      `json:"tenant_id" mapstructure:"tenant_id"`
}

type RouterPage struct {
	pagination.LinkedPageBase
}

func (p RouterPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"routers_links"`
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

func (p RouterPage) IsEmpty() (bool, error) {
	is, err := ExtractRouters(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractRouters(page pagination.Page) ([]Router, error) {
	var resp struct {
		Routers []Router `mapstructure:"routers" json:"routers"`
	}

	err := mapstructure.Decode(page.(RouterPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Routers, nil
}

type commonResult struct {
	gophercloud.CommonResult
}

// Extract is a function that accepts a result and extracts a network resource.
func (r commonResult) Extract() (*Router, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Router *Router `json:"router"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Neutron router: %v", err)
	}

	return res.Router, nil
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
