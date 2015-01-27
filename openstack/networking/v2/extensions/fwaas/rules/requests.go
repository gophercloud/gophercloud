package rules

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type ListOpts struct {
	TenantID             string `q:"tenant_id"`
	Name                 string `q:"name"`
	Description          string `q:"description"`
	Protocol             string `q:"protocol"`
	Action               string `q:"action"`
	IpVersion            int    `q:"ip_version"`
	SourceIpAddress      string `q:"source_ip_address"`
	DestinationIpAddress string `q:"destination_ip_address"`
	SourcePort           string `q:"source_port"`
	DestinationPort      string `q:"destination_port"`
	Enabled              bool   `q:"enabled"`
	ID                   string `q:"id"`
	Limit                int    `q:"limit"`
	Marker               string `q:"marker"`
	SortKey              string `q:"sort_key"`
	SortDir              string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// firewall rules. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
//
// Default policy settings return only those firewall rules that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return RulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOpts contains all the values needed to create a new firewall rule.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a firewall rule
	// for another tenant.
	TenantId             string
	Name                 string
	Description          string
	Protocol             string
	Action               string
	IpVersion            int
	SourceIpAddress      string
	DestinationIpAddress string
	SourcePort           string
	DestinationPort      string
	Shared               bool
	Enabled              bool
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall rule
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type rule struct {
		TenantId             string `json:"tenant_id,omitempty"`
		Name                 string `json:"name,omitempty"`
		Description          string `json:"description,omitempty"`
		Protocol             string `json:"protocol"`
		Action               string `json:"action"`
		IpVersion            int    `json:"ip_version,omitempty"`
		SourceIpAddress      string `json:"source_ip_address,omitempty"`
		DestinationIpAddress string `json:"destination_ip_address,omitempty"`
		SourcePort           string `json:"source_port,omitempty"`
		DestinationPort      string `json:"destination_port,omitempty"`
		Shared               bool   `json:"shared,omitempty"`
		Enabled              bool   `json:"enabled,omitempty"`
	}
	type request struct {
		Rule rule `json:"firewall_rule"`
	}

	reqBody := request{Rule: rule{
		TenantId:             opts.TenantId,
		Name:                 opts.Name,
		Description:          opts.Description,
		Protocol:             opts.Protocol,
		Action:               opts.Action,
		IpVersion:            opts.IpVersion,
		SourceIpAddress:      opts.SourceIpAddress,
		DestinationIpAddress: opts.DestinationIpAddress,
		SourcePort:           opts.SourcePort,
		DestinationPort:      opts.DestinationPort,
		Shared:               opts.Shared,
		Enabled:              opts.Enabled,
	}}

	var res CreateResult
	_, res.Err = perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{201},
	})
	return res
}

// Get retrieves a particular firewall rule based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOpts contains the values used when updating a firewall rule.
type UpdateOpts struct {
	Name                 string
	Description          string
	Protocol             string
	Action               string
	IpVersion            int
	SourceIpAddress      string
	DestinationIpAddress string
	SourcePort           string
	DestinationPort      string
	Shared               bool
	Enabled              bool
}

// Update allows firewall policies to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type rule struct {
		Name                 string `json:"name"`
		Description          string `json:"description"`
		Protocol             string `json:"protocol,omitempty"`
		Action               string `json:"action,omitempty"`
		IpVersion            int    `json:"ip_version,omitempty"`
		SourceIpAddress      string `json:"source_ip_address,omitempty"`
		DestinationIpAddress string `json:"destination_ip_address,omitempty"`
		SourcePort           string `json:"source_port,omitempty"`
		DestinationPort      string `json:"destination_port,omitempty"`
		Shared               bool   `json:"shared,omitempty"`
		Enabled              bool   `json:"enabled,omitempty"`
	}
	type request struct {
		Rule rule `json:"firewall_rule"`
	}

	reqBody := request{Rule: rule{
		Name:                 opts.Name,
		Description:          opts.Description,
		Protocol:             opts.Protocol,
		Action:               opts.Action,
		IpVersion:            opts.IpVersion,
		SourceIpAddress:      opts.SourceIpAddress,
		DestinationIpAddress: opts.DestinationIpAddress,
		SourcePort:           opts.SourcePort,
		DestinationPort:      opts.DestinationPort,
		Shared:               opts.Shared,
		Enabled:              opts.Enabled,
	}}

	// Send request to API
	var res UpdateResult
	_, res.Err = perigee.Request("PUT", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

// Delete will permanently delete a particular firewall rule based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
