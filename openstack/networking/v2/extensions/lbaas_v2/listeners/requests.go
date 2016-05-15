package listeners

import (
	"fmt"

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
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID              string `q:"id"`
	Name            string `q:"name"`
	AdminStateUp    *bool  `q:"admin_state_up"`
	TenantID        string `q:"tenant_id"`
	DefaultPoolID   string `q:"default_pool_id"`
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
		return ListenerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

var (
	errLoadbalancerIdRequired = fmt.Errorf("Loadbalancer ID is required")
	errProtocolRequired       = fmt.Errorf("Protocol is required")
	errProtocolPortRequired   = fmt.Errorf("Protocol port is required")
)

// CreateOpts contains all the values needed to create a new Listener.
type CreateOpts struct {
	// Required. The protocol - can either be TCP, HTTP or HTTPS.
	Protocol string

	// Required. The port on which to listen for client traffic.
	ProtocolPort int

	// Required for admins. Indicates the owner of the Listener.
	TenantID string

	// Required. The load balancer on which to provision this listener.
	LoadbalancerID string

	// Human-readable name for the Listener. Does not have to be unique.
	Name string

	// Optional. The ID of the default pool with which the Listener is associated.
	DefaultPoolID string

	// Optional. Human-readable description for the Listener.
	Description string

	// Optional. The maximum number of connections allowed for the Listener.
	ConnLimit *int

	// Optional. A reference to a container of TLS secrets.
	DefaultTlsContainerRef string

	// Optional. A list of references to TLS secrets.
	SniContainerRefs []string

	// Optional. The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Create is an operation which provisions a new Listeners based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create Listeners on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	// Validate required opts
	if opts.LoadbalancerID == "" {
		res.Err = errLoadbalancerIdRequired
		return res
	}
	if opts.Protocol == "" {
		res.Err = errProtocolRequired
		return res
	}
	if opts.ProtocolPort == 0 {
		res.Err = errProtocolPortRequired
		return res
	}

	type listener struct {
		Name                   *string  `json:"name,omitempty"`
		LoadbalancerID         string   `json:"loadbalancer_id,omitempty"`
		Protocol               string   `json:"protocol"`
		ProtocolPort           int      `json:"protocol_port"`
		DefaultPoolID          *string  `json:"default_pool_id,omitempty"`
		Description            *string  `json:"description,omitempty"`
		TenantID               *string  `json:"tenant_id,omitempty"`
		ConnLimit              *int     `json:"connection_limit,omitempty"`
		AdminStateUp           *bool    `json:"admin_state_up,omitempty"`
		DefaultTlsContainerRef *string  `json:"default_tls_container_ref,omitempty"`
		SniContainerRefs       []string `json:"sni_container_refs,omitempty"`
	}

	type request struct {
		Listener listener `json:"listener"`
	}

	reqBody := request{Listener: listener{
		Name:                   gophercloud.MaybeString(opts.Name),
		LoadbalancerID:         opts.LoadbalancerID,
		Protocol:               opts.Protocol,
		ProtocolPort:           opts.ProtocolPort,
		DefaultPoolID:          gophercloud.MaybeString(opts.DefaultPoolID),
		Description:            gophercloud.MaybeString(opts.Description),
		TenantID:               gophercloud.MaybeString(opts.TenantID),
		ConnLimit:              opts.ConnLimit,
		AdminStateUp:           opts.AdminStateUp,
		DefaultTlsContainerRef: gophercloud.MaybeString(opts.DefaultTlsContainerRef),
		SniContainerRefs:       opts.SniContainerRefs,
	}}

	_, res.Err = c.Post(rootURL(c), reqBody, &res.Body, nil)
	return res
}

// Get retrieves a particular Listeners based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, id), &res.Body, nil)
	return res
}

// UpdateOpts contains all the values that can be updated on an existing Listener.
// Attributes not listed here but appear in CreateOpts are immutable and cannot
// be updated.
type UpdateOpts struct {
	// Human-readable name for the Listener. Does not have to be unique.
	Name string

	// Optional. Human-readable description for the Listener.
	Description string

	// Optional. The maximum number of connections allowed for the Listener.
	ConnLimit *int

	// Optional. A reference to a container of TLS secrets.
	DefaultTlsContainerRef string

	// Optional. A list of references to TLS secrets.
	SniContainerRefs []string

	// Optional. The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Update is an operation which modifies the attributes of the specified Listener.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type listener struct {
		Name                   string   `json:"name,omitempty"`
		Description            *string  `json:"description,omitempty"`
		ConnLimit              *int     `json:"connection_limit,omitempty"`
		AdminStateUp           *bool    `json:"admin_state_up,omitempty"`
		DefaultTlsContainerRef *string  `json:"default_tls_container_ref,omitempty"`
		SniContainerRefs       []string `json:"sni_container_refs,omitempty"`
	}

	type request struct {
		Listener listener `json:"listener"`
	}

	reqBody := request{Listener: listener{
		Name:                   opts.Name,
		Description:            gophercloud.MaybeString(opts.Description),
		ConnLimit:              opts.ConnLimit,
		AdminStateUp:           opts.AdminStateUp,
		DefaultTlsContainerRef: gophercloud.MaybeString(opts.DefaultTlsContainerRef),
		SniContainerRefs:       opts.SniContainerRefs,
	}}

	var res UpdateResult
	_, res.Err = c.Put(resourceURL(c, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})

	return res
}

// Delete will permanently delete a particular Listeners based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, id), nil)
	return res
}
