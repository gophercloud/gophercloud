package subnets

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the subnet attributes you want to see returned. SortKey allows you to sort
// by a particular subnet attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Name       string `q:"name"`
	EnableDHCP *bool  `q:"enable_dhcp"`
	NetworkID  string `q:"network_id"`
	TenantID   string `q:"tenant_id"`
	IPVersion  int    `q:"ip_version"`
	GatewayIP  string `q:"gateway_ip"`
	CIDR       string `q:"cidr"`
	ID         string `q:"id"`
	Limit      int    `q:"limit"`
	Marker     string `q:"marker"`
	SortKey    string `q:"sort_key"`
	SortDir    string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// subnets. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those subnets that are owned by the tenant
// who submits the request, unless the request is submitted by an user with
// administrative rights.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	// Build query parameters

	u := listURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return SubnetPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

// Get retrieves a specific subnet based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	return res
}

// Valid IP types
const (
	IPv4 = 4
	IPv6 = 6
)

// CreateOpts represents the attributes used when creating a new subnet.
type CreateOpts struct {
	// Required
	NetworkID string
	CIDR      string
	// Optional
	Name            string
	TenantID        string
	AllocationPools []AllocationPool
	GatewayIP       string
	IPVersion       int
	EnableDHCP      *bool
	DNSNameservers  []string
	HostRoutes      []interface{}
}

// Create accepts a CreateOpts struct and creates a new subnet using the values
// provided. You must remember to provide a valid NetworkID, CIDR and IP version.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	// Validate required options
	if opts.NetworkID == "" {
		res.Err = errNetworkIDRequired
		return res
	}
	if opts.CIDR == "" {
		res.Err = errCIDRRequired
		return res
	}
	if opts.IPVersion != 0 && opts.IPVersion != IPv4 && opts.IPVersion != IPv6 {
		res.Err = errInvalidIPType
		return res
	}

	type subnet struct {
		NetworkID       string           `json:"network_id"`
		CIDR            string           `json:"cidr"`
		Name            *string          `json:"name,omitempty"`
		TenantID        *string          `json:"tenant_id,omitempty"`
		AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
		GatewayIP       *string          `json:"gateway_ip,omitempty"`
		IPVersion       int              `json:"ip_version,omitempty"`
		EnableDHCP      *bool            `json:"enable_dhcp,omitempty"`
		DNSNameservers  []string         `json:"dns_nameservers,omitempty"`
		HostRoutes      []interface{}    `json:"host_routes,omitempty"`
	}
	type request struct {
		Subnet subnet `json:"subnet"`
	}

	reqBody := request{Subnet: subnet{
		NetworkID:  opts.NetworkID,
		CIDR:       opts.CIDR,
		Name:       gophercloud.MaybeString(opts.Name),
		TenantID:   gophercloud.MaybeString(opts.TenantID),
		GatewayIP:  gophercloud.MaybeString(opts.GatewayIP),
		EnableDHCP: opts.EnableDHCP,
	}}

	if opts.IPVersion != 0 {
		reqBody.Subnet.IPVersion = opts.IPVersion
	}
	if len(opts.AllocationPools) != 0 {
		reqBody.Subnet.AllocationPools = opts.AllocationPools
	}
	if len(opts.DNSNameservers) != 0 {
		reqBody.Subnet.DNSNameservers = opts.DNSNameservers
	}
	if len(opts.HostRoutes) != 0 {
		reqBody.Subnet.HostRoutes = opts.HostRoutes
	}

	_, res.Err = perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})

	return res
}

// UpdateOpts represents the attributes used when updating an existing subnet.
type UpdateOpts struct {
	Name           string
	GatewayIP      string
	DNSNameservers []string
	HostRoutes     []interface{}
	EnableDHCP     *bool
}

// Update accepts a UpdateOpts struct and updates an existing subnet using the
// values provided.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type subnet struct {
		Name           *string       `json:"name,omitempty"`
		GatewayIP      *string       `json:"gateway_ip,omitempty"`
		DNSNameservers []string      `json:"dns_nameservers,omitempty"`
		HostRoutes     []interface{} `json:"host_routes,omitempty"`
		EnableDHCP     *bool         `json:"enable_dhcp,omitempty"`
	}
	type request struct {
		Subnet subnet `json:"subnet"`
	}

	reqBody := request{Subnet: subnet{
		Name:       gophercloud.MaybeString(opts.Name),
		GatewayIP:  gophercloud.MaybeString(opts.GatewayIP),
		EnableDHCP: opts.EnableDHCP,
	}}

	if len(opts.DNSNameservers) != 0 {
		reqBody.Subnet.DNSNameservers = opts.DNSNameservers
	}

	if len(opts.HostRoutes) != 0 {
		reqBody.Subnet.HostRoutes = opts.HostRoutes
	}

	var res UpdateResult
	_, res.Err = perigee.Request("PUT", updateURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200, 201},
	})

	return res
}

// Delete accepts a unique ID and deletes the subnet associated with it.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
