package networks

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type ListOpts struct {
	Status       string
	Name         string
	AdminStateUp *bool
	TenantID     string
	Shared       *bool
	ID           string
	Page         int
	PerPage      int
}

type NetworkOpts struct {
	AdminStateUp bool
	Name         string
	Shared       *bool
	TenantID     string
}

func ptrToStr(val *bool) string {
	if *val == true {
		return "true"
	} else if *val == false {
		return "false"
	} else {
		return ""
	}
}

func List(c *gophercloud.ServiceClient, opts ListOpts) gophercloud.Pager {
	// Build query parameters
	q := make(map[string]string)
	if opts.Status != "" {
		q["status"] = opts.Status
	}
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = ptrToStr(opts.AdminStateUp)
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.Shared != nil {
		q["shared"] = ptrToStr(opts.Shared)
	}
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.Page != 0 {
		q["page"] = strconv.Itoa(opts.Page)
	}
	if opts.PerPage != 0 {
		q["per_page"] = strconv.Itoa(opts.PerPage)
	}

	u := ListURL(c) + utils.BuildQuery(q)
	return gophercloud.NewLinkedPager(c, u)
}

func Get(c *gophercloud.ServiceClient, id string) (*Network, error) {
	var n Network
	_, err := perigee.Request("GET", GetURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Network *Network `json:"network"`
		}{&n},
		OkCodes: []int{200},
	})
	if err != nil {
		return nil, err
	}
	return &n, nil
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

func Update(c *gophercloud.ServiceClient, networkID string, opts NetworkOpts) (*Network, error) {
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
		Network *Network `json:"network"`
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
	_, err := perigee.Request("PUT", GetURL(c, networkID), perigee.Options{
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

func Delete(c *gophercloud.ServiceClient, networkID string) error {
	_, err := perigee.Request("DELETE", DeleteURL(c, networkID), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
