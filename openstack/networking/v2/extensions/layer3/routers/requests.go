package routers

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

const (
	version      = "v2.0"
	resourcePath = "routers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id)
}

type ListOpts struct {
	ID           string
	Name         string
	AdminStateUp *bool
	Status       string
	TenantID     string
	Limit        int
	Marker       string
	SortKey      string
	SortDir      string
}

func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q := make(map[string]string)
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = strconv.FormatBool(*opts.AdminStateUp)
	}
	if opts.Status != "" {
		q["status"] = opts.Status
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.Marker != "" {
		q["marker"] = opts.Marker
	}
	if opts.Limit != 0 {
		q["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.SortKey != "" {
		q["sort_key"] = opts.SortKey
	}
	if opts.SortDir != "" {
		q["sort_dir"] = opts.SortDir
	}

	u := rootURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return RouterPage{pagination.LinkedPageBase(r)}
	})
}

type CreateOpts struct {
	Name         string
	AdminStateUp *bool
	TenantID     string
	GatewayInfo  *GatewayInfo
}

func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type router struct {
		Name         *string      `json:"name,omitempty"`
		AdminStateUp *bool        `json:"admin_state_up,omitempty"`
		TenantID     *string      `json:"tenant_id,omitempty"`
		GatewayInfo  *GatewayInfo `json:"external_gateway_info,omitempty"`
	}

	type request struct {
		Router router `json:"router"`
	}

	reqBody := request{Router: router{
		Name:         gophercloud.MaybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		TenantID:     gophercloud.MaybeString(opts.TenantID),
	}}

	if opts.GatewayInfo != nil {
		reqBody.Router.GatewayInfo = opts.GatewayInfo
	}

	var res CreateResult
	_, err := perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})
	res.Err = err
	return res
}

func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, err := perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	res.Err = err
	return res
}

type UpdateOpts struct {
	Name         string
	AdminStateUp *bool
	GatewayInfo  *GatewayInfo
}

func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type router struct {
		Name         *string      `json:"name,omitempty"`
		AdminStateUp *bool        `json:"admin_state_up,omitempty"`
		GatewayInfo  *GatewayInfo `json:"external_gateway_info,omitempty"`
	}

	type request struct {
		Router router `json:"router"`
	}

	reqBody := request{Router: router{
		Name:         gophercloud.MaybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
	}}

	if opts.GatewayInfo != nil {
		reqBody.Router.GatewayInfo = opts.GatewayInfo
	}

	// Send request to API
	var res UpdateResult
	_, err := perigee.Request("PUT", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	res.Err = err
	return res
}

func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, err := perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Err = err
	return res
}
