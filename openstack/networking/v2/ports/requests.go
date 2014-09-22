package ports

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type ListOpts struct {
	Status          string
	Name            string
	AdminStateUp    *bool
	NetworkID       string
	TenantID        string
	DeviceOwner     string
	MACAddress      string
	ID              string
	SecurityGroups  string
	DeviceID        string
	BindingHostID   string
	BindingVIFType  string
	BindingVNICType string
	Limit           int
	Page            string
	PerPage         string
	SortKey         string
	SortDir         string
}

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
	if opts.SecurityGroups != "" {
		q["security_groups"] = opts.SecurityGroups
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
	if opts.Page != "" {
		q["page"] = opts.Page
	}
	if opts.PerPage != "" {
		q["per_page"] = opts.PerPage
	}
	if opts.SortKey != "" {
		q["sort_key"] = opts.SortKey
	}
	if opts.SortDir != "" {
		q["sort_dir"] = opts.SortDir
	}

	u := listURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return PortPage{pagination.LinkedPageBase(r)}
	})
}

func Get(c *gophercloud.ServiceClient, id string) (*Port, error) {
	var p Port
	_, err := perigee.Request("GET", getURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Port *Port `json:"port"`
		}{&p},
		OkCodes: []int{200},
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func maybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

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

func Create(c *gophercloud.ServiceClient, opts CreateOpts) (*Port, error) {
	type port struct {
		NetworkID      string      `json:"network_id,omitempty"`
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
		return nil, ErrNetworkIDRequired
	}

	// Populate request body
	reqBody := request{Port: port{
		NetworkID:    opts.NetworkID,
		Name:         maybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		TenantID:     maybeString(opts.TenantID),
		MACAddress:   maybeString(opts.MACAddress),
		DeviceID:     maybeString(opts.DeviceID),
		DeviceOwner:  maybeString(opts.DeviceOwner),
	}}

	if opts.FixedIPs != nil {
		reqBody.Port.FixedIPs = opts.FixedIPs
	}

	if opts.SecurityGroups != nil {
		reqBody.Port.SecurityGroups = opts.SecurityGroups
	}

	// Response
	type response struct {
		Port *Port `json:"port"`
	}
	var res response
	_, err := perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res,
		OkCodes:     []int{201},
		DumpReqJson: true,
	})
	if err != nil {
		return nil, err
	}

	return res.Port, nil
}

type UpdateOpts struct {
	Name           string
	AdminStateUp   *bool
	FixedIPs       interface{}
	DeviceID       string
	DeviceOwner    string
	SecurityGroups []string
}

func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) (*Port, error) {
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
		Name:         maybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		DeviceID:     maybeString(opts.DeviceID),
		DeviceOwner:  maybeString(opts.DeviceOwner),
	}}

	if opts.FixedIPs != nil {
		reqBody.Port.FixedIPs = opts.FixedIPs
	}

	if opts.SecurityGroups != nil {
		reqBody.Port.SecurityGroups = opts.SecurityGroups
	}

	// Response
	type response struct {
		Port *Port `json:"port"`
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

	return res.Port, nil
}

func Delete(c *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
