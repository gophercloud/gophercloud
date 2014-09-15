package networks

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

func APIVersions(c *gophercloud.ServiceClient) (*APIVersionsList, error) {
	var resp APIVersionsList
	_, err := perigee.Request("GET", APIVersionsURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func APIInfo(c *gophercloud.ServiceClient, v string) (*APIInfoList, error) {
	var resp APIInfoList
	_, err := perigee.Request("GET", APIInfoURL(c, v), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetExtension(c *gophercloud.ServiceClient, name string) (*Extension, error) {
	var ext Extension
	_, err := perigee.Request("GET", ExtensionURL(c, name), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Extension *Extension `json:"extension"`
		}{&ext},
		OkCodes: []int{200},
	})

	if err != nil {
		return nil, err
	}
	return &ext, nil
}

func Get(c *gophercloud.ServiceClient, id string) (*NetworkResult, error) {
	var n NetworkResult
	_, err := perigee.Request("GET", NetworkURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Network *NetworkResult `json:"network"`
		}{&n},
		OkCodes: []int{200},
	})
	if err != nil {
		return nil, err
	}
	return &n, nil
}

type NetworkOpts struct {
	AdminStateUp bool
	Name         string
	Shared       *bool
	TenantID     string
}

func Create(c *gophercloud.ServiceClient, opts NetworkOpts) (*NetworkCreateResult, error) {
	// Define structures
	type network struct {
		AdminStateUp bool    `json:"admin_state_up"`
		Name         string  `json:"name"`
		Shared       *bool   `json:"shared,omitempty"`
		TenantID     *string `json:"tenant_id,omitempty"`
	}
	type request struct {
		Network network `json:"network"`
	}
	type response struct {
		Network *NetworkCreateResult `json:"network"`
	}

	// Validate
	if opts.Name == "" {
		return nil, ErrNameRequired
	}

	// Populate request body
	reqBody := request{Network: network{
		AdminStateUp: opts.AdminStateUp,
		Name:         opts.Name,
		Shared:       opts.Shared,
	}}

	if opts.TenantID != "" {
		reqBody.Network.TenantID = &opts.TenantID
	}

	// Send request to API
	var res response
	_, err := perigee.Request("POST", CreateURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res,
		OkCodes:     []int{201},
	})
	if err != nil {
		return nil, err
	}

	return res.Network, nil
}

func Update(c *gophercloud.ServiceClient, networkID string, opts NetworkOpts) (*NetworkResult, error) {
	// Define structures
	type network struct {
		AdminStateUp bool    `json:"admin_state_up"`
		Name         string  `json:"name"`
		Shared       *bool   `json:"shared,omitempty"`
		TenantID     *string `json:"tenant_id,omitempty"`
	}

	type request struct {
		Network network `json:"network"`
	}
	type response struct {
		Network *NetworkResult `json:"network"`
	}

	// Populate request body
	reqBody := request{Network: network{
		AdminStateUp: opts.AdminStateUp,
		Name:         opts.Name,
		Shared:       opts.Shared,
	}}

	if opts.TenantID != "" {
		reqBody.Network.TenantID = &opts.TenantID
	}

	// Send request to API
	var res response
	_, err := perigee.Request("PUT", NetworkURL(c, networkID), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return res.Network, nil
}
