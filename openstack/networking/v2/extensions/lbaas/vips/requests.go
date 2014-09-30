package vips

import (
	"strconv"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID              string
	Name            string
	AdminStateUp    *bool
	Status          string
	TenantID        string
	SubnetID        string
	Address         string
	PortID          string
	Protocol        string
	ProtocolPort    int
	ConnectionLimit int
	Limit           int
	Marker          string
	SortKey         string
	SortDir         string
}

// List returns a Pager which allows you to iterate over a collection of
// routers. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those routers that are owned by the
// tenant who submits the request, unless an admin user submits the request.
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
	if opts.SubnetID != "" {
		q["subnet_id"] = opts.SubnetID
	}
	if opts.Address != "" {
		q["address"] = opts.Address
	}
	if opts.PortID != "" {
		q["port_id"] = opts.PortID
	}
	if opts.Protocol != "" {
		q["protocol"] = opts.Protocol
	}
	if opts.ProtocolPort != 0 {
		q["protocol_port"] = strconv.Itoa(opts.ProtocolPort)
	}
	if opts.ConnectionLimit != 0 {
		q["connection_limit"] = strconv.Itoa(opts.ConnectionLimit)
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
		return VIPPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}
