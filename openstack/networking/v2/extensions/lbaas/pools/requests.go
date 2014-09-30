package pools

import (
	"strconv"

	"github.com/racker/perigee"
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
	Status       string
	LBMethod     string
	Protocol     string
	SubnetID     string
	TenantID     string
	AdminStateUp *bool
	Name         string
	ID           string
	VIPID        string
	Limit        int
	Marker       string
	SortKey      string
	SortDir      string
}

// List returns a Pager which allows you to iterate over a collection of
// pools. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those pools that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q := make(map[string]string)

	if opts.Status != "" {
		q["status"] = opts.Status
	}
	if opts.LBMethod != "" {
		q["lb_method"] = opts.LBMethod
	}
	if opts.Protocol != "" {
		q["protocol"] = opts.Protocol
	}
	if opts.SubnetID != "" {
		q["subnet_id"] = opts.SubnetID
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = strconv.FormatBool(*opts.AdminStateUp)
	}
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.VIPID != "" {
		q["vip_id"] = opts.VIPID
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
		return PoolPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

const (
	LBMethodRoundRobin       = "ROUND_ROBIN"
	LBMethodLeastConnections = "LEAST_CONNECTIONS"

	ProtocolTCP   = "TCP"
	ProtocolHTTP  = "HTTP"
	ProtocolHTTPS = "HTTPS"
)

// CreateOpts contains all the values needed to create a new pool.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a pool
	// for another tenant.
	TenantID string

	// Required. Name of the pool.
	Name string

	// Required. The protocol used by the pool members, you can use either
	// ProtocolTCP, ProtocolHTTP, or ProtocolHTTPS.
	Protocol string

	// The network on which the members of the pool will be located. Only members
	// that are on this network can be added to the pool.
	SubnetID string

	// The algorithm used to distribute load between the members of the pool. The
	// current specification supports LBMethodRoundRobin and
	// LBMethodLeastConnections as valid values for this attribute.
	LBMethod string
}

// Create accepts a CreateOpts struct and uses the values to create a new
// load balancer pool.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type pool struct {
		Name     string `json:"name,"`
		TenantID string `json:"tenant_id"`
		Protocol string `json:"protocol"`
		SubnetID string `json:"subnet_id"`
		LBMethod string `json:"lb_method"`
	}
	type request struct {
		Pool pool `json:"pool"`
	}

	reqBody := request{Pool: pool{
		Name:     opts.Name,
		TenantID: opts.TenantID,
		Protocol: opts.Protocol,
		SubnetID: opts.SubnetID,
		LBMethod: opts.LBMethod,
	}}

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
