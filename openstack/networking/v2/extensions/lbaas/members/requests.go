package members

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
	Weight       int
	AdminStateUp *bool
	TenantID     string
	PoolID       string
	Address      string
	ProtocolPort int
	ID           string
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
	if opts.Weight != 0 {
		q["weight"] = strconv.Itoa(opts.Weight)
	}
	if opts.PoolID != "" {
		q["pool_id"] = opts.PoolID
	}
	if opts.Address != "" {
		q["address"] = opts.Address
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = strconv.FormatBool(*opts.AdminStateUp)
	}
	if opts.ProtocolPort != 0 {
		q["protocol_port"] = strconv.Itoa(opts.ProtocolPort)
	}
	if opts.ID != "" {
		q["id"] = opts.ID
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
		return MemberPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

// CreateOpts contains all the values needed to create a new pool member.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a pool
	// for another tenant.
	TenantID string

	// Required. The IP address of the member.
	Address string

	// Required. The port on which the application is hosted.
	ProtocolPort int

	// Required. The pool to which this member will belong.
	PoolID string
}

// Create accepts a CreateOpts struct and uses the values to create a new
// load balancer pool member.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type member struct {
		TenantID     string `json:"tenant_id"`
		ProtocolPort int    `json:"protocol_port"`
		Address      string `json:"address"`
		PoolID       string `json:"pool_id"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{
		Address:      opts.Address,
		TenantID:     opts.TenantID,
		ProtocolPort: opts.ProtocolPort,
		PoolID:       opts.PoolID,
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

// Get retrieves a particular pool member based on its unique ID.
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

// UpdateOpts contains the values used when updating a pool member.
type UpdateOpts struct {
	// The administrative state of the member, which is up (true) or down (false).
	AdminStateUp bool
}

// Update allows members to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type member struct {
		AdminStateUp bool `json:"admin_state_up"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{AdminStateUp: opts.AdminStateUp}}

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

// Delete will permanently delete a particular member based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, err := perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Err = err
	return res
}
