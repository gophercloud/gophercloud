package policies

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type ListOpts struct {
	TenantID    string `q:"tenant_id"`
	Name        string `q:"name"`
	Description string `q:"description"`
	Shared      bool   `q:"shared"`
	Audited     bool   `q:"audited"`
	ID          string `q:"id"`
	Limit       int    `q:"limit"`
	Marker      string `q:"marker"`
	SortKey     string `q:"sort_key"`
	SortDir     string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// firewall policies. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
//
// Default policy settings return only those firewall policies that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOpts contains all the values needed to create a new firewall policy.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a firewall policy
	// for another tenant.
	TenantID    string
	Name        string
	Description string
	Shared      *bool
	Audited     *bool
	Rules       []string
}

// ToPolicyCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToPolicyCreateMap() map[string]interface{} {
	p := make(map[string]interface{})

	if opts.TenantID != "" {
		p["tenant_id"] = opts.TenantID
	}
	if opts.Name != "" {
		p["name"] = opts.Name
	}
	if opts.Description != "" {
		p["description"] = opts.Description
	}
	if opts.Shared != nil {
		p["shared"] = *opts.Shared
	}
	if opts.Audited != nil {
		p["audited"] = *opts.Audited
	}
	if opts.Rules != nil {
		p["firewall_rules"] = opts.Rules
	}

	return p
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall policy
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {

	type request struct {
		Policy map[string]interface{} `json:"firewall_policy"`
	}

	reqBody := request{Policy: opts.ToPolicyCreateMap()}

	var res CreateResult
	_, res.Err = perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{201},
	})
	return res
}

// Get retrieves a particular firewall policy based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOpts contains the values used when updating a firewall policy.
type UpdateOpts struct {
	// Name of the firewall policy.
	Name        *string
	Description *string
	Shared      *bool
	Audited     *bool
	Rules       []string
}

// ToPolicyUpdateMap casts a CreateOpts struct to a map.
func (opts UpdateOpts) ToPolicyUpdateMap() map[string]interface{} {
	p := make(map[string]interface{})

	if opts.Name != nil {
		p["name"] = *opts.Name
	}
	if opts.Description != nil {
		p["description"] = *opts.Description
	}
	if opts.Shared != nil {
		p["shared"] = *opts.Shared
	}
	if opts.Audited != nil {
		p["audited"] = *opts.Audited
	}
	if opts.Rules != nil {
		p["firewall_rules"] = opts.Rules
	}

	return p
}

// Update allows firewall policies to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {

	type request struct {
		Policy map[string]interface{} `json:"firewall_policy"`
	}

	reqBody := request{Policy: opts.ToPolicyUpdateMap()}

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

// Delete will permanently delete a particular firewall policy based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}

func InsertRule(c *gophercloud.ServiceClient, policyID, ruleID, beforeID, afterID string) error {
	type request struct {
		RuleId string `json:"firewall_rule_id"`
		Before string `json:"insert_before,omitempty"`
		After  string `json:"insert_after,omitempty"`
	}

	reqBody := request{
		RuleId: ruleID,
		Before: beforeID,
		After:  afterID,
	}

	// Send request to API
	var res commonResult
	_, res.Err = perigee.Request("PUT", insertURL(c, policyID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res.Err
}

func RemoveRule(c *gophercloud.ServiceClient, policyID, ruleID string) error {
	type request struct {
		RuleId string `json:"firewall_rule_id"`
	}

	reqBody := request{
		RuleId: ruleID,
	}

	// Send request to API
	var res commonResult
	_, res.Err = perigee.Request("PUT", removeURL(c, policyID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res.Err
}
