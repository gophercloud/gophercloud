package networks

import (
	"fmt"
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type networkOpts struct {
	AdminStateUp *bool
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

type CreateOptsInt interface {
	ToMap() map[string]map[string]interface{}
	IsCreateOpts() bool
}

type CreateOpts networkOpts

func (o CreateOpts) ToMap() map[string]map[string]interface{} {
	inner := make(map[string]interface{})

	if o.AdminStateUp != nil {
		inner["admin_state_up"] = &o.AdminStateUp
	}
	if o.Name != "" {
		inner["name"] = o.Name
	}
	if o.Shared != nil {
		inner["shared"] = &o.Shared
	}
	if o.TenantID != "" {
		inner["tenant_id"] = o.TenantID
	}

	outer := make(map[string]map[string]interface{})
	outer["network"] = inner

	return outer
}

func (o CreateOpts) IsCreateOpts() bool { return true }

// Create accepts a CreateOpts struct and creates a new network using the values
// provided. This operation does not actually require a request body, i.e. the
// CreateOpts struct argument can be empty.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// network. An admin user, however, has the option of specifying another tenant
// ID in the CreateOpts struct.
func Create(c *gophercloud.ServiceClient, opts CreateOptsInt) CreateResult {
	var res CreateResult

	if opts.IsCreateOpts() != true {
		res.Err = fmt.Errorf("Must provide valid create opts")
		return res
	}

	reqBody := opts.ToMap()

	// Send request to API
	_, err := perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})
	res.Err = err
	return res
}

type UpdateOptsInt interface {
	ToMap() map[string]map[string]interface{}
	IsUpdateOpts() bool
}

type UpdateOpts networkOpts

func (o UpdateOpts) ToMap() map[string]map[string]interface{} {
	inner := make(map[string]interface{})

	if o.AdminStateUp != nil {
		inner["admin_state_up"] = &o.AdminStateUp
	}
	if o.Name != "" {
		inner["name"] = o.Name
	}
	if o.Shared != nil {
		inner["shared"] = &o.Shared
	}

	outer := make(map[string]map[string]interface{})
	outer["network"] = inner

	return outer
}

func (o UpdateOpts) IsUpdateOpts() bool { return true }

// Update accepts a UpdateOpts struct and updates an existing network using the
// values provided. For more information, see the Create function.
func Update(c *gophercloud.ServiceClient, networkID string, opts UpdateOptsInt) UpdateResult {
	var res UpdateResult

	if opts.IsUpdateOpts() != true {
		res.Err = fmt.Errorf("Must provide valid update opts")
		return res
	}

	reqBody := opts.ToMap()

	// Send request to API
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
