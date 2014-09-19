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
