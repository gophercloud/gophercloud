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
