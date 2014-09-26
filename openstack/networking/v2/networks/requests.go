package networks

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type networkOpts struct {
	AdminStateUp bool
	Name         string
	Shared       *bool
	TenantID     string
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the network attributes you want to see returned. SortKey allows you to sort
// by a particular network attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status       string
	Name         string
	AdminStateUp *bool
	TenantID     string
	Shared       *bool
	ID           string
	Marker       string
	Limit        int
	SortKey      string
	SortDir      string
}

// List returns a Pager which allows you to iterate over a collection of
// networks. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	// Build query parameters
	q := make(map[string]string)
	if opts.Status != "" {
		q["status"] = opts.Status
	}
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = strconv.FormatBool(*opts.AdminStateUp)
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.Shared != nil {
		q["shared"] = strconv.FormatBool(*opts.Shared)
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

	u := listURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

// Get retrieves a specific network based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, err := perigee.Request("GET", getURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	res.Err = err
	return res
}

type CreateOpts interface {
	ToMap() map[string]interface{}
}

// CreateOpts represents the attributes used when creating a new network.
type CreateOpts networkOpts

// Create accepts a CreateOpts struct and creates a new network using the values
// provided. This operation does not actually require a request body, i.e. the
// CreateOpts struct argument can be empty.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// network. An admin user, however, has the option of specifying another tenant
// ID in the CreateOpts struct.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	// Define structures
	type network struct {
		AdminStateUp bool    `json:"admin_state_up,omitempty"`
		Name         string  `json:"name,omitempty"`
		Shared       *bool   `json:"shared,omitempty"`
		TenantID     *string `json:"tenant_id,omitempty"`
	}
	type request struct {
		Network network `json:"network"`
	}

	// Populate request body
	reqBody := request{Network: network{
		AdminStateUp: opts.AdminStateUp,
		Name:         opts.Name,
		Shared:       opts.Shared,
	}}

	if opts.TenantID != "" {
		reqBody.Network.TenantID = &opts.TenantID
	}

	// Send request to API
	var res CreateResult
	_, err := perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})
	res.Err = err
	return res
}

// UpdateOpts represents the attributes used when updating an existing network.
type UpdateOpts networkOpts

// Update accepts a UpdateOpts struct and updates an existing network using the
// values provided. For more information, see the Create function.
func Update(c *gophercloud.ServiceClient, networkID string, opts UpdateOpts) UpdateResult {
	// Define structures
	type network struct {
		AdminStateUp bool   `json:"admin_state_up"`
		Name         string `json:"name"`
		Shared       *bool  `json:"shared,omitempty"`
	}

	type request struct {
		Network network `json:"network"`
	}

	// Populate request body
	reqBody := request{Network: network{
		AdminStateUp: opts.AdminStateUp,
		Name:         opts.Name,
		Shared:       opts.Shared,
	}}

	// Send request to API
	var res UpdateResult
	_, err := perigee.Request("PUT", getURL(c, networkID), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200, 201},
	})
	res.Err = err
	return res
}

// Delete accepts a unique ID and deletes the network associated with it.
func Delete(c *gophercloud.ServiceClient, networkID string) DeleteResult {
	var res DeleteResult
	_, err := perigee.Request("DELETE", deleteURL(c, networkID), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Err = err
	return res
}
