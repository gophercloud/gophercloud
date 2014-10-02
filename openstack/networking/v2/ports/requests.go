package ports

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the port attributes you want to see returned. SortKey allows you to sort
// by a particular port attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status          string
	Name            string
	AdminStateUp    *bool
	NetworkID       string
	TenantID        string
	DeviceOwner     string
	MACAddress      string
	ID              string
	DeviceID        string
	BindingHostID   string
	BindingVIFType  string
	BindingVNICType string
	Limit           int
	Marker          string
	SortKey         string
	SortDir         string
}

// List returns a Pager which allows you to iterate over a collection of
// ports. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those ports that are owned by the tenant
// who submits the request, unless the request is submitted by an user with
// administrative rights.
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
	if opts.NetworkID != "" {
		q["network_id"] = opts.NetworkID
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
	}
	if opts.DeviceOwner != "" {
		q["device_owner"] = opts.DeviceOwner
	}
	if opts.MACAddress != "" {
		q["mac_address"] = opts.MACAddress
	}
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.DeviceID != "" {
		q["device_id"] = opts.DeviceID
	}
	if opts.BindingHostID != "" {
		q["binding:host_id"] = opts.BindingHostID
	}
	if opts.BindingVIFType != "" {
		q["binding:vif_type"] = opts.BindingVIFType
	}
	if opts.BindingVNICType != "" {
		q["binding:vnic_type"] = opts.BindingVNICType
	}
	if opts.NetworkID != "" {
		q["network_id"] = opts.NetworkID
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
		return PortPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

// Get retrieves a specific port based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	return res
}

// CreateOpts represents the attributes used when creating a new port.
type CreateOpts struct {
	NetworkID      string
	Name           string
	AdminStateUp   *bool
	MACAddress     string
	FixedIPs       interface{}
	DeviceID       string
	DeviceOwner    string
	TenantID       string
	SecurityGroups []string
}

// Create accepts a CreateOpts struct and creates a new network using the values
// provided. You must remember to provide a NetworkID value.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	type port struct {
		NetworkID      string      `json:"network_id"`
		Name           *string     `json:"name,omitempty"`
		AdminStateUp   *bool       `json:"admin_state_up,omitempty"`
		MACAddress     *string     `json:"mac_address,omitempty"`
		FixedIPs       interface{} `json:"fixed_ips,omitempty"`
		DeviceID       *string     `json:"device_id,omitempty"`
		DeviceOwner    *string     `json:"device_owner,omitempty"`
		TenantID       *string     `json:"tenant_id,omitempty"`
		SecurityGroups []string    `json:"security_groups,omitempty"`
	}
	type request struct {
		Port port `json:"port"`
	}

	// Validate
	if opts.NetworkID == "" {
		res.Err = errNetworkIDRequired
		return res
	}

	// Populate request body
	reqBody := request{Port: port{
		NetworkID:    opts.NetworkID,
		Name:         gophercloud.MaybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		TenantID:     gophercloud.MaybeString(opts.TenantID),
		MACAddress:   gophercloud.MaybeString(opts.MACAddress),
		DeviceID:     gophercloud.MaybeString(opts.DeviceID),
		DeviceOwner:  gophercloud.MaybeString(opts.DeviceOwner),
	}}

	if opts.FixedIPs != nil {
		reqBody.Port.FixedIPs = opts.FixedIPs
	}

	if opts.SecurityGroups != nil {
		reqBody.Port.SecurityGroups = opts.SecurityGroups
	}

	// Response
	_, res.Err = perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
		DumpReqJson: true,
	})

	return res
}

// UpdateOpts represents the attributes used when updating an existing port.
type UpdateOpts struct {
	Name           string
	AdminStateUp   *bool
	FixedIPs       interface{}
	DeviceID       string
	DeviceOwner    string
	SecurityGroups []string
}

// Update accepts a UpdateOpts struct and updates an existing port using the
// values provided.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type port struct {
		Name           *string     `json:"name,omitempty"`
		AdminStateUp   *bool       `json:"admin_state_up,omitempty"`
		FixedIPs       interface{} `json:"fixed_ips,omitempty"`
		DeviceID       *string     `json:"device_id,omitempty"`
		DeviceOwner    *string     `json:"device_owner,omitempty"`
		SecurityGroups []string    `json:"security_groups,omitempty"`
	}
	type request struct {
		Port port `json:"port"`
	}

	// Populate request body
	reqBody := request{Port: port{
		Name:         gophercloud.MaybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		DeviceID:     gophercloud.MaybeString(opts.DeviceID),
		DeviceOwner:  gophercloud.MaybeString(opts.DeviceOwner),
	}}

	if opts.FixedIPs != nil {
		reqBody.Port.FixedIPs = opts.FixedIPs
	}

	if opts.SecurityGroups != nil {
		reqBody.Port.SecurityGroups = opts.SecurityGroups
	}

	// Response
	var res UpdateResult
	_, res.Err = perigee.Request("PUT", updateURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200, 201},
	})
	return res
}

// Delete accepts a unique ID and deletes the port associated with it.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
