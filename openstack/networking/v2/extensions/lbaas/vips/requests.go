package vips

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// AdminState gives users a solid type to work with for create and update
// operations. It is recommended that users use the `Up` and `Down` enums.
type AdminState *bool

// Convenience vars for AdminStateUp values.
var (
	iTrue  = true
	iFalse = false

	Up   AdminState = &iTrue
	Down AdminState = &iFalse
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID              string `q:"id"`
	Name            string `q:"name"`
	AdminStateUp    *bool  `q:"admin_state_up"`
	Status          string `q:"status"`
	TenantID        string `q:"tenant_id"`
	SubnetID        string `q:"subnet_id"`
	Address         string `q:"address"`
	PortID          string `q:"port_id"`
	Protocol        string `q:"protocol"`
	ProtocolPort    int    `q:"protocol_port"`
	ConnectionLimit int    `q:"connection_limit"`
	Limit           int    `q:"limit"`
	Marker          string `q:"marker"`
	SortKey         string `q:"sort_key"`
	SortDir         string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// routers. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those routers that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return VIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is what types must satisfy to be used as Create
// options.
type CreateOptsBuilder interface {
	ToVIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new virtual IP.
type CreateOpts struct {
	// Required. Human-readable name for the VIP. Does not have to be unique.
	Name string

	// Required. The network on which to allocate the VIP's address. A tenant can
	// only create VIPs on networks authorized by policy (e.g. networks that
	// belong to them or networks that are shared).
	SubnetID string

	// Required. The protocol - can either be TCP, HTTP or HTTPS.
	Protocol string

	// Required. The port on which to listen for client traffic.
	ProtocolPort int

	// Required. The ID of the pool with which the VIP is associated.
	PoolID string

	// Required for admins. Indicates the owner of the VIP.
	TenantID string

	// Optional. The IP address of the VIP.
	Address string

	// Optional. Human-readable description for the VIP.
	Description string

	// Optional. Omit this field to prevent session persistence.
	Persistence *SessionPersistence

	// Optional. The maximum number of connections allowed for the VIP.
	ConnLimit *int

	// Optional. The administrative state of the VIP. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// ToVIPCreateMap allows CreateOpts to satisfy the CreateOptsBuilder
// interface
func (opts CreateOpts) ToVIPCreateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		err := gophercloud.ErrMissingInput{}
		err.Function = "vips.ToVIPCreateMap"
		err.Argument = "vips.CreateOpts.Name"
		return nil, err
	}
	if opts.SubnetID == "" {
		err := gophercloud.ErrMissingInput{}
		err.Function = "vips.ToVIPCreateMap"
		err.Argument = "vips.CreateOpts.SubnetID"
		return nil, err
	}
	if opts.Protocol == "" {
		err := gophercloud.ErrMissingInput{}
		err.Function = "vips.ToVIPCreateMap"
		err.Argument = "vips.CreateOpts.Protocol"
		return nil, err
	}
	if opts.ProtocolPort == 0 {
		err := gophercloud.ErrMissingInput{}
		err.Function = "vips.ToVIPCreateMap"
		err.Argument = "vips.CreateOpts.ProtocolPort"
		return nil, err
	}
	if opts.PoolID == "" {
		err := gophercloud.ErrMissingInput{}
		err.Function = "vips.ToVIPCreateMap"
		err.Argument = "vips.CreateOpts.PoolID"
		return nil, err
	}

	i := map[string]interface{}{
		"name":          opts.Name,
		"subnet_id":     opts.SubnetID,
		"protocol":      opts.Protocol,
		"protocol_port": opts.ProtocolPort,
		"pool_id":       opts.PoolID,
	}
	if opts.Description != "" {
		i["description"] = opts.Description
	}
	if opts.TenantID != "" {
		i["tenant_id"] = opts.TenantID
	}
	if opts.Address != "" {
		i["address"] = opts.Address
	}
	if opts.Persistence != nil {
		i["session_persistence"] = opts.Persistence
	}
	if opts.ConnLimit != nil {
		i["connection_limit"] = opts.ConnLimit
	}
	if opts.AdminStateUp != nil {
		i["admin_state_up"] = opts.AdminStateUp
	}

	b := make(map[string]interface{})
	b["vip"] = i

	return b, nil
}

// Create is an operation which provisions a new virtual IP based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Please note that the PoolID should refer to a pool that is not already
// associated with another vip. If the pool is already used by another vip,
// then the operation will fail with a 409 Conflict error will be returned.
//
// Users with an admin role can create VIPs on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var r CreateResult

	b, err := opts.ToVIPCreateMap()
	if err != nil {
		r.Err = err
		return r
	}

	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return r
}

// Get retrieves a particular virtual IP based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, id), &res.Body, nil)
	return res
}

// UpdateOpts contains all the values needed to update an existing virtual IP.
// Attributes not listed here but appear in CreateOpts are immutable and cannot
// be updated.
type UpdateOpts struct {
	// Human-readable name for the VIP. Does not have to be unique.
	Name string

	// Required. The ID of the pool with which the VIP is associated.
	PoolID string

	// Optional. Human-readable description for the VIP.
	Description string

	// Optional. Omit this field to prevent session persistence.
	Persistence *SessionPersistence

	// Optional. The maximum number of connections allowed for the VIP.
	ConnLimit *int

	// Optional. The administrative state of the VIP. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Update is an operation which modifies the attributes of the specified VIP.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type vip struct {
		Name         string              `json:"name,omitempty"`
		PoolID       string              `json:"pool_id,omitempty"`
		Description  *string             `json:"description,omitempty"`
		Persistence  *SessionPersistence `json:"session_persistence,omitempty"`
		ConnLimit    *int                `json:"connection_limit,omitempty"`
		AdminStateUp *bool               `json:"admin_state_up,omitempty"`
	}

	type request struct {
		VirtualIP vip `json:"vip"`
	}

	reqBody := request{VirtualIP: vip{
		Name:         opts.Name,
		PoolID:       opts.PoolID,
		Description:  gophercloud.MaybeString(opts.Description),
		ConnLimit:    opts.ConnLimit,
		AdminStateUp: opts.AdminStateUp,
	}}

	if opts.Persistence != nil {
		reqBody.VirtualIP.Persistence = opts.Persistence
	}

	var res UpdateResult
	_, res.Err = c.Put(resourceURL(c, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})

	return res
}

// Delete will permanently delete a particular virtual IP based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, id), nil)
	return res
}
