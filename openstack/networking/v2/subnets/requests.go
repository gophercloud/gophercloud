package subnets

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type ListOpts struct {
	Name       string
	EnableDHCP *bool
	NetworkID  string
	TenantID   string
	IPVersion  int
	GatewayIP  string
	CIDR       string
	ID         string
	Limit      int
	Page       string
	PerPage    string
	SortKey    string
	SortDir    string
}

func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	// Build query parameters
	q := make(map[string]string)
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.EnableDHCP != nil {
		q["enable_dhcp"] = strconv.FormatBool(*opts.EnableDHCP)
	}
	if opts.NetworkID != "" {
		q["network_id"] = opts.NetworkID
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.IPVersion != 0 {
		q["ip_version"] = strconv.Itoa(opts.IPVersion)
	}
	if opts.GatewayIP != "" {
		q["gateway_ip"] = opts.GatewayIP
	}
	if opts.CIDR != "" {
		q["cidr"] = opts.CIDR
	}
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.Limit != 0 {
		q["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.Page != "" {
		q["page"] = opts.Page
	}
	if opts.PerPage != "" {
		q["per_page"] = opts.PerPage
	}
	if opts.SortKey != "" {
		q["sort_key"] = opts.SortKey
	}
	if opts.SortDir != "" {
		q["sort_dir"] = opts.SortDir
	}

	u := ListURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return SubnetPage{pagination.LinkedPageBase(r)}
	})
}

func Get(c *gophercloud.ServiceClient, id string) (*Subnet, error) {
	var s Subnet
	_, err := perigee.Request("GET", GetURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Subnet *Subnet `json:"subnet"`
		}{&s},
		OkCodes: []int{200},
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

const (
	IPv4 = 4
	IPv6 = 6
)

type SubnetOpts struct {
	// Required
	NetworkID string
	CIDR      string
	// Optional
	Name            string
	TenantID        string
	AllocationPools []AllocationPool
	GatewayIP       string
	IPVersion       int
	ID              string
	EnableDHCP      *bool
}

// maybeString returns nil for empty strings and nil for empty.
func maybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

func Create(c *gophercloud.ServiceClient, opts SubnetOpts) (*Subnet, error) {
	// Validate required options
	if opts.NetworkID == "" {
		return nil, ErrNetworkIDRequired
	}
	if opts.CIDR == "" {
		return nil, ErrCIDRRequired
	}
	if opts.IPVersion != 0 && opts.IPVersion != IPv4 && opts.IPVersion != IPv6 {
		return nil, ErrInvalidIPType
	}

	type subnet struct {
		NetworkID       string           `json:"network_id"`
		CIDR            string           `json:"cidr"`
		Name            *string          `json:"name,omitempty"`
		TenantID        *string          `json:"tenant_id,omitempty"`
		AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
		GatewayIP       *string          `json:"gateway_ip,omitempty"`
		IPVersion       int              `json:"ip_version,omitempty"`
		ID              *string          `json:"id,omitempty"`
		EnableDHCP      *bool            `json:"enable_dhcp,omitempty"`
	}
	type request struct {
		Subnet subnet `json:"subnet"`
	}

	reqBody := request{Subnet: subnet{
		NetworkID:  opts.NetworkID,
		CIDR:       opts.CIDR,
		Name:       maybeString(opts.Name),
		TenantID:   maybeString(opts.TenantID),
		GatewayIP:  maybeString(opts.GatewayIP),
		ID:         maybeString(opts.ID),
		EnableDHCP: opts.EnableDHCP,
	}}

	if opts.IPVersion != 0 {
		reqBody.Subnet.IPVersion = opts.IPVersion
	}

	if len(opts.AllocationPools) != 0 {
		reqBody.Subnet.AllocationPools = opts.AllocationPools
	}

	type response struct {
		Subnet *Subnet `json:"subnet"`
	}

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

	return res.Subnet, nil
}

func Update(c *gophercloud.ServiceClient, id string, opts SubnetOpts) (*Subnet, error) {
	if opts.CIDR != "" {
		return nil, ErrCIDRNotUpdatable
	}
	if opts.IPVersion != 0 {
		return nil, ErrIPVersionNotUpdatable
	}

	type subnet struct {
		NetworkID       string           `json:"network_id,omitempty"`
		Name            *string          `json:"name,omitempty"`
		TenantID        *string          `json:"tenant_id,omitempty"`
		AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
		GatewayIP       *string          `json:"gateway_ip,omitempty"`
		ID              *string          `json:"id,omitempty"`
		EnableDHCP      *bool            `json:"enable_dhcp,omitempty"`
	}
	type request struct {
		Subnet subnet `json:"subnet"`
	}

	reqBody := request{Subnet: subnet{
		NetworkID:  opts.NetworkID,
		Name:       maybeString(opts.Name),
		TenantID:   maybeString(opts.TenantID),
		GatewayIP:  maybeString(opts.GatewayIP),
		ID:         maybeString(opts.ID),
		EnableDHCP: opts.EnableDHCP,
	}}

	if len(opts.AllocationPools) != 0 {
		reqBody.Subnet.AllocationPools = opts.AllocationPools
	}

	type response struct {
		Subnet *Subnet `json:"subnet"`
	}

	var res response
	_, err := perigee.Request("PUT", UpdateURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res,
		OkCodes:     []int{200, 201},
	})
	if err != nil {
		return nil, err
	}

	return res.Subnet, nil
}

func Delete(c *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", DeleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
