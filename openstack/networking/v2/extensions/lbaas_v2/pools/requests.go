package pools

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
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
// the Pool attributes you want to see returned. SortKey allows you to
// sort by a particular Pool attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	LBMethod       string `q:"lb_algorithm"`
	Protocol       string `q:"protocol"`
	SubnetID       string `q:"subnet_id"`
	TenantID       string `q:"tenant_id"`
	AdminStateUp   *bool  `q:"admin_state_up"`
	Name           string `q:"name"`
	ID             string `q:"id"`
	LoadbalancerID string `q:"loadbalancer_id"`
	Limit          int    `q:"limit"`
	Marker         string `q:"marker"`
	SortKey        string `q:"sort_key"`
	SortDir        string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// pools. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those pools that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Supported attributes for create/update operations.
const (
	LBMethodRoundRobin       = "ROUND_ROBIN"
	LBMethodLeastConnections = "LEAST_CONNECTIONS"
	LBMethodSourceIp         = "SOURCE_IP"

	ProtocolTCP   = "TCP"
	ProtocolHTTP  = "HTTP"
	ProtocolHTTPS = "HTTPS"
)

// CreateOpts contains all the values needed to create a new pool.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a pool
	// for another tenant.
	TenantID string

	// Optional. The network on which the members of the pool will be located.
	// Only members that are on this network can be added to the pool.
	SubnetID string

	// Optional. Name of the pool.
	Name string

	// Optional. Human-readable description for the pool.
	Description string

	// Required. The protocol used by the pool members, you can use either
	// ProtocolTCP, ProtocolHTTP, or ProtocolHTTPS.
	Protocol string

	// The Loadbalancer on which the members of the pool will be associated with.
	// Note:  one of LoadbalancerID or ListenerID must be provided.
	LoadbalancerID string

	// The Listener on which the members of the pool will be associated with.
	// Note:  one of LoadbalancerID or ListenerID must be provided.
	ListenerID string

	// The algorithm used to distribute load between the members of the pool. The
	// current specification supports LBMethodRoundRobin, LBMethodLeastConnections
	// and LBMethodSourceIp as valid values for this attribute.
	LBMethod string

	// Optional. Omit this field to prevent session persistence.
	Persistence *SessionPersistence

	// Optional. The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Create accepts a CreateOpts struct and uses the values to create a new
// load balancer pool.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type pool struct {
		Name           string              `json:"name,omitempty"`
		Description    string              `json:"description,omitempty"`
		TenantID       string              `json:"tenant_id,omitempty"`
		SubnetID       string              `json:"subnet_id,omitempty"`
		Protocol       string              `json:"protocol"`
		LoadbalancerID string              `json:"loadbalancer_id,omitempty"`
		ListenerID     string              `json:"listener_id,omitempty"`
		LBMethod       string              `json:"lb_algorithm"`
		Persistence    *SessionPersistence `json:"session_persistence,omitempty"`
		AdminStateUp   *bool               `json:"admin_state_up,omitempty"`
	}

	type request struct {
		Pool pool `json:"pool"`
	}

	reqBody := request{Pool: pool{
		Name:           opts.Name,
		Description:    opts.Description,
		TenantID:       opts.TenantID,
		SubnetID:       opts.SubnetID,
		Protocol:       opts.Protocol,
		LoadbalancerID: opts.LoadbalancerID,
		ListenerID:     opts.ListenerID,
		LBMethod:       opts.LBMethod,
		AdminStateUp:   opts.AdminStateUp,
	}}

	if opts.Persistence != nil {
		reqBody.Pool.Persistence = opts.Persistence
	}

	var res CreateResult
	_, res.Err = c.Post(rootURL(c), reqBody, &res.Body, nil)
	return res
}

// Get retrieves a particular pool based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, id), &res.Body, nil)
	return res
}

// UpdateOpts contains the values used when updating a pool.
type UpdateOpts struct {
	// Optional. Name of the pool.
	Name string

	// Optional. Human-readable description for the pool.
	Description string

	// The algorithm used to distribute load between the members of the pool. The
	// current specification supports LBMethodRoundRobin, LBMethodLeastConnections
	// and LBMethodSourceIp as valid values for this attribute.
	LBMethod string

	// Optional. The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Update allows pools to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type pool struct {
		Name         string `json:"name,omitempty"`
		Description  string `json:"description,omitempty"`
		LBMethod     string `json:"lb_algorithm,omitempty"`
		AdminStateUp *bool  `json:"admin_state_up,omitempty"`
	}
	type request struct {
		Pool pool `json:"pool"`
	}

	reqBody := request{Pool: pool{
		Name:         opts.Name,
		Description:  opts.Description,
		LBMethod:     opts.LBMethod,
		AdminStateUp: opts.AdminStateUp,
	}}

	// Send request to API
	var res UpdateResult
	_, res.Err = c.Put(resourceURL(c, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Delete will permanently delete a particular pool based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, id), nil)
	return res
}

// MemberListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Member attributes you want to see returned. SortKey allows you to
// sort by a particular Member attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type MemberListOpts struct {
	Name         string `q:"name"`
	Weight       int    `q:"weight"`
	AdminStateUp *bool  `q:"admin_state_up"`
	TenantID     string `q:"tenant_id"`
	SubnetID     string `q:"subnet_id"`
	Address      string `q:"address"`
	ProtocolPort int    `q:"protocol_port"`
	ID           string `q:"id"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// members. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those members that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func ListAssociateMembers(c *gophercloud.ServiceClient, poolID string, opts MemberListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := memberRootURL(c, poolID) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return MemberPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOpts contains all the values needed to create a new Member for a Pool.
type MemberCreateOpts struct {
	// Optional. Name of the Member.
	Name string

	// Only required if the caller has an admin role and wants to create a Member
	// for another tenant.
	TenantID string

	// Required. The IP address of the member to receive traffic from the load balancer.
	Address string

	// Required. The port on which to listen for client traffic.
	ProtocolPort int

	// Optional. A positive integer value that indicates the relative portion of
	// traffic that this member should receive from the pool. For example, a
	// member with a weight of 10 receives five times as much traffic as a member
	// with a weight of 2.
	Weight int

	// Optional.  If you omit this parameter, LBaaS uses the vip_subnet_id
	// parameter value for the subnet UUID.
	SubnetID string

	// Optional. The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// CreateAssociateMember will create and associate a Member with a particular Pool.
func CreateAssociateMember(c *gophercloud.ServiceClient, poolID string, opts MemberCreateOpts) AssociateResult {
	type member struct {
		Name         string `json:"name,omitempty"`
		TenantID     string `json:"tenant_id,omitempty"`
		Address      string `json:"address,omitempty"`
		ProtocolPort int    `json:"protocol_port,omitempty"`
		Weight       int    `json:"weight,omitempty"`
		SubnetID     string `json:"subnet_id,omitempty"`
		AdminStateUp *bool  `json:"admin_state_up,omitempty"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{
		Name:         opts.Name,
		TenantID:     opts.TenantID,
		Address:      opts.Address,
		ProtocolPort: opts.ProtocolPort,
		Weight:       opts.Weight,
		SubnetID:     opts.SubnetID,
		AdminStateUp: opts.AdminStateUp,
	}}

	var res AssociateResult
	_, res.Err = c.Post(memberRootURL(c, poolID), reqBody, &res.Body, nil)
	return res
}

// Get retrieves a particular Pool Member based on its unique ID.
func GetAssociateMember(c *gophercloud.ServiceClient, poolID string, memberID string) GetResult {
	var res GetResult
	_, res.Err = c.Get(memberResourceURL(c, poolID, memberID), &res.Body, nil)
	return res
}

// UpdateOpts contains the values used when updating a Pool Member.
type MemberUpdateOpts struct {
	// Name of the Member.
	Name string

	// A positive integer value that indicates the relative portion of
	// traffic that this member should receive from the pool. For example, a
	// member with a weight of 10 receives five times as much traffic as a member
	// with a weight of 2.
	Weight int

	// The administrative state of the member, which is up (true) or down (false).
	AdminStateUp *bool
}

// Update allows Member to be updated.
func UpdateAssociateMember(c *gophercloud.ServiceClient, poolID string, memberID string, opts MemberUpdateOpts) UpdateResult {
	type member struct {
		Name         string `json:"name,omitempty"`
		Weight       int    `json:"weight,omitempty"`
		AdminStateUp *bool  `json:"admin_state_up,omitempty"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{
		Name:         opts.Name,
		Weight:       opts.Weight,
		AdminStateUp: opts.AdminStateUp,
	}}

	// Send request to API
	var res UpdateResult
	_, res.Err = c.Put(memberResourceURL(c, poolID, memberID), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return res
}

// DisassociateMember will remove and disassociate a Member from a particular Pool.
func DeleteMember(c *gophercloud.ServiceClient, poolID string, memberID string) AssociateResult {
	var res AssociateResult
	_, res.Err = c.Delete(memberResourceURL(c, poolID, memberID), nil)
	return res
}
