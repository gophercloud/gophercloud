package subnets

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

func maybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the subnet attributes you want to see returned. SortKey allows you to sort
// by a particular subnet attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
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
	Marker     string
	SortKey    string
	SortDir    string
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
	if opts.Marker != "" {
		q["marker"] = opts.Marker
	}
	if opts.SortKey != "" {
		q["sort_key"] = opts.SortKey
	}
	if opts.SortDir != "" {
		q["sort_dir"] = opts.SortDir
	}

	u := listURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return SubnetPage{pagination.LinkedPageBase(r)}
	})
}

// Get retrieves a specific subnet based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (*Subnet, error) {
	var s Subnet
	_, err := perigee.Request("GET", getURL(c, id), perigee.Options{
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
func Create(c *gophercloud.ServiceClient, opts CreateOpts) (*Subnet, error) {
	// Validate required options
	if opts.NetworkID == "" {
		return nil, errNetworkIDRequired
	}
	if opts.CIDR == "" {
		return nil, errCIDRRequired
	}
	if opts.IPVersion != 0 && opts.IPVersion != IPv4 && opts.IPVersion != IPv6 {
		return nil, errInvalidIPType
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
		Name:       maybeString(opts.Name),
		TenantID:   maybeString(opts.TenantID),
		GatewayIP:  maybeString(opts.GatewayIP),
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

	type response struct {
		Subnet *Subnet `json:"subnet"`
	}

	var res response
	_, err := perigee.Request("POST", createURL(c), perigee.Options{
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
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) (*Subnet, error) {
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
		Name:       maybeString(opts.Name),
		GatewayIP:  maybeString(opts.GatewayIP),
		EnableDHCP: opts.EnableDHCP,
	}}

	if len(opts.DNSNameservers) != 0 {
		reqBody.Subnet.DNSNameservers = opts.DNSNameservers
	}

	if len(opts.HostRoutes) != 0 {
		reqBody.Subnet.HostRoutes = opts.HostRoutes
	}

	type response struct {
		Subnet *Subnet `json:"subnet"`
	}

	var res response
	_, err := perigee.Request("PUT", updateURL(c, id), perigee.Options{
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

// Delete accepts a unique ID and deletes the subnet associated with it.
func Delete(c *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
