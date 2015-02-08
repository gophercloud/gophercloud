package rules

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Firewall rule attributes you want to see returned. SortKey allows you to
// sort by a particular firewall rule attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	TenantID             string `q:"tenant_id"`
	Name                 string `q:"name"`
	Description          string `q:"description"`
	Protocol             string `q:"protocol"`
	Action               string `q:"action"`
	IPVersion            int    `q:"ip_version"`
	SourceIPAddress      string `q:"source_ip_address"`
	DestinationIPAddress string `q:"destination_ip_address"`
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
	// Mandatory for create
	Protocol string
	Action   string
	// Optional
	TenantID             string
	Name                 string
	Description          string
	IPVersion            int
	SourceIPAddress      string
	DestinationIPAddress string
	SourcePort           string
	DestinationPort      string
	Shared               *bool
	Enabled              *bool
}

// ToRuleCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToRuleCreateMap() map[string]interface{} {
	r := make(map[string]interface{})

	r["protocol"] = opts.Protocol
	r["action"] = opts.Action

	if opts.TenantID != "" {
		r["tenant_id"] = opts.TenantID
	}
	if opts.Name != "" {
		r["name"] = opts.Name
	}
	if opts.Description != "" {
		r["description"] = opts.Description
	}
	if opts.IPVersion != 0 {
		r["ip_version"] = opts.IPVersion
	}
	if opts.SourceIPAddress != "" {
		r["source_ip_address"] = opts.SourceIPAddress
	}
	if opts.DestinationIPAddress != "" {
		r["destination_ip_address"] = opts.DestinationIPAddress
	}
	if opts.SourcePort != "" {
		r["source_port"] = opts.SourcePort
	}
	if opts.DestinationPort != "" {
		r["destination_port"] = opts.DestinationPort
	}
	if opts.Shared != nil {
		r["shared"] = *opts.Shared
	}
	if opts.Enabled != nil {
		r["enabled"] = *opts.Enabled
	}

	return r
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall rule
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {

	type request struct {
		Rule map[string]interface{} `json:"firewall_rule"`
	}

	reqBody := request{Rule: opts.ToRuleCreateMap()}

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
// Optional
type UpdateOpts struct {
	Protocol             string
	Action               string
	Name                 *string
	Description          *string
	IPVersion            int
	SourceIPAddress      *string
	DestinationIPAddress *string
	SourcePort           *string
	DestinationPort      *string
	Shared               *bool
	Enabled              *bool
}

// ToRuleUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToRuleUpdateMap() map[string]interface{} {
	r := make(map[string]interface{})

	if opts.Protocol != "" {
		r["protocol"] = opts.Protocol
	}
	if opts.Action != "" {
		r["action"] = opts.Action
	}
	if opts.Name != nil {
		r["name"] = *opts.Name
	}
	if opts.Description != nil {
		r["description"] = *opts.Description
	}
	if opts.IPVersion != 0 {
		r["ip_version"] = opts.IPVersion
	}
	if opts.SourceIPAddress != nil {
		s := *opts.SourceIPAddress
		if s == "" {
			r["source_ip_address"] = nil
		} else {
			r["source_ip_address"] = s
		}
	}
	if opts.DestinationIPAddress != nil {
		s := *opts.DestinationIPAddress
		if s == "" {
			r["destination_ip_address"] = nil
		} else {
			r["destination_ip_address"] = s
		}
	}
	if opts.SourcePort != nil {
		s := *opts.SourcePort
		if s == "" {
			r["source_port"] = nil
		} else {
			r["source_port"] = s
		}
	}
	if opts.DestinationPort != nil {
		s := *opts.DestinationPort
		if s == "" {
			r["destination_port"] = nil
		} else {
			r["destination_port"] = s
		}
	}
	if opts.Shared != nil {
		r["shared"] = *opts.Shared
	}
	if opts.Enabled != nil {
		r["enabled"] = *opts.Enabled
	}

	return r
}

// Update allows firewall policies to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {

	type request struct {
		Rule map[string]interface{} `json:"firewall_rule"`
	}

	reqBody := request{Rule: opts.ToRuleUpdateMap()}

	// Send request to API
	var res UpdateResult
	_, res.Err = perigee.Request("PUT", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
		DumpReqJson: true,
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
