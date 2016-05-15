package loadbalancers

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
// the Loadbalancer attributes you want to see returned. SortKey allows you to
// sort by a particular attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Description        string `q:"description"`
	AdminStateUp       *bool  `q:"admin_state_up"`
	TenantID           string `q:"tenant_id"`
	ProvisioningStatus string `q:"provisioning_status"`
	VipAddress         string `q:"vip_address"`
	VipSubnetID        string `q:"vip_subnet_id"`
	ID                 string `q:"id"`
	OperatingStatus    string `q:"operating_status"`
	Name               string `q:"name"`
	Flavor             string `q:"flavor"`
	Provider           string `q:"provider"`
	Limit              int    `q:"limit"`
	Marker             string `q:"marker"`
	SortKey            string `q:"sort_key"`
	SortDir            string `q:"sort_dir"`
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
		return LoadbalancerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

var (
	errVipSubnetIDRequried = fmt.Errorf("VipSubnetID is required")
)

// CreateOpts contains all the values needed to create a new Loadbalancer.
type CreateOpts struct {
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string

	// Optional. Human-readable description for the Loadbalancer.
	Description string

	// Required. The network on which to allocate the Loadbalancer's address. A tenant can
	// only create Loadbalancers on networks authorized by policy (e.g. networks that
	// belong to them or networks that are shared).
	VipSubnetID string

	// Required for admins. The UUID of the tenant who owns the Loadbalancer.
	// Only administrative users can specify a tenant UUID other than their own.
	TenantID string

	// Optional. The IP address of the Loadbalancer.
	VipAddress string

	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool

	// Optional. The UUID of a flavor.
	Flavor string

	// Optional. The name of the provider.
	Provider string
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	// Validate required opts
	if opts.VipSubnetID == "" {
		res.Err = errVipSubnetIDRequried
		return res
	}

	type loadbalancer struct {
		Name         *string `json:"name,omitempty"`
		Description  *string `json:"description,omitempty"`
		VipSubnetID  string  `json:"vip_subnet_id"`
		TenantID     *string `json:"tenant_id,omitempty"`
		VipAddress   *string `json:"vip_address,omitempty"`
		AdminStateUp *bool   `json:"admin_state_up,omitempty"`
		Flavor       *string `json:"flavor,omitempty"`
		Provider     *string `json:"provider,omitempty"`
	}

	type request struct {
		Loadbalancer loadbalancer `json:"loadbalancer"`
	}

	reqBody := request{Loadbalancer: loadbalancer{
		Name:         gophercloud.MaybeString(opts.Name),
		Description:  gophercloud.MaybeString(opts.Description),
		VipSubnetID:  opts.VipSubnetID,
		TenantID:     gophercloud.MaybeString(opts.TenantID),
		VipAddress:   gophercloud.MaybeString(opts.VipAddress),
		AdminStateUp: opts.AdminStateUp,
		Flavor:       gophercloud.MaybeString(opts.Flavor),
		Provider:     gophercloud.MaybeString(opts.Provider),
	}}

	_, res.Err = c.Post(rootURL(c), reqBody, &res.Body, nil)
	return res
}

// Get retrieves a particular virtual IP based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, id), &res.Body, nil)
	return res
}

// UpdateOpts contains all the values needed to update an existing virtual
// Loadbalancer. Attributes not listed here but appear in CreateOpts are
// immutable and cannot be updated.
type UpdateOpts struct {
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string

	// Optional. Human-readable description for the Loadbalancer.
	Description string

	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool
}

// Update is an operation which modifies the attributes of the specified Loadbalancer.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type loadbalancer struct {
		Name         *string `json:"name,omitempty"`
		Description  *string `json:"description,omitempty"`
		AdminStateUp *bool   `json:"admin_state_up,omitempty"`
	}

	type request struct {
		Loadbalancer loadbalancer `json:"loadbalancer"`
	}

	reqBody := request{Loadbalancer: loadbalancer{
		Name:         gophercloud.MaybeString(opts.Name),
		Description:  gophercloud.MaybeString(opts.Description),
		AdminStateUp: opts.AdminStateUp,
	}}

	var res UpdateResult
	_, res.Err = c.Put(resourceURL(c, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})

	return res
}

// Delete will permanently delete a particular Loadbalancer based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, id), nil)
	return res
}
