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

	u := ListURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return PortPage{pagination.LinkedPageBase(r)}
	})
}

func Get(c *gophercloud.ServiceClient, id string) (*Port, error) {
	var p Port
	_, err := perigee.Request("GET", GetURL(c, id), perigee.Options{
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

type PortOpts struct {
	NetworkID      string
	Status         string
	Name           string
	AdminStateUp   *bool
	TenantID       string
	MACAddress     string
	FixedIPs       interface{}
	SecurityGroups []string
}

func maybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

func Create(c *gophercloud.ServiceClient, opts PortOpts) (*Port, error) {
	type port struct {
		NetworkID      string      `json:"network_id,omitempty"`
		Status         *string     `json:"status,omitempty"`
		Name           *string     `json:"name,omitempty"`
		AdminStateUp   *bool       `json:"admin_state_up,omitempty"`
		TenantID       *string     `json:"tenant_id,omitempty"`
		MACAddress     *string     `json:"mac_address,omitempty"`
		FixedIPs       interface{} `json:"fixed_ips,omitempty"`
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
		Status:       maybeString(opts.Status),
		Name:         maybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		TenantID:     maybeString(opts.TenantID),
		MACAddress:   maybeString(opts.MACAddress),
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
	_, err := perigee.Request("POST", CreateURL(c), perigee.Options{
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

func Update(c *gophercloud.ServiceClient, id string, opts PortOpts) (*Port, error) {
	type port struct {
		NetworkID      string      `json:"network_id,omitempty"`
		Status         *string     `json:"status,omitempty"`
		Name           *string     `json:"name,omitempty"`
		AdminStateUp   *bool       `json:"admin_state_up,omitempty"`
		TenantID       *string     `json:"tenant_id,omitempty"`
		MACAddress     *string     `json:"mac_address,omitempty"`
		FixedIPs       interface{} `json:"fixed_ips,omitempty"`
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
		Status:       maybeString(opts.Status),
		Name:         maybeString(opts.Name),
		AdminStateUp: opts.AdminStateUp,
		TenantID:     maybeString(opts.TenantID),
		MACAddress:   maybeString(opts.MACAddress),
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
	_, err := perigee.Request("PUT", UpdateURL(c, id), perigee.Options{
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
	_, err := perigee.Request("DELETE", DeleteURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
