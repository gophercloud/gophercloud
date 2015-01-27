package firewalls

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type ListOpts struct {
	TenantID     string `q:"tenant_id"`
	Name         string `q:"name"`
	Description  string `q:"description"`
	AdminStateUp bool   `q:"admin_state_up"`
	Shared       bool   `q:"shared"`
	PolicyId     string `q:"policy_id"`
	ID           string `q:"id"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// firewalls. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
//
// Default policy settings return only those firewalls that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return FirewallPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOpts contains all the values needed to create a new firewall.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a firewall
	// for another tenant.
	TenantId     string
	Name         string
	Description  string
	AdminStateUp bool
	Shared       bool
	PolicyId     string
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type firewall struct {
		TenantId     string `json:"tenant_id,omitempty"`
		Name         string `json:"name,omitempty"`
		Description  string `json:"description,omitempty"`
		AdminStateUp bool   `json:"admin_state_up,omitempty"`
		PolicyId     string `json:"firewall_policy_id"`
	}
	type request struct {
		Firewall firewall `json:"firewall"`
	}

	reqBody := request{Firewall: firewall{
		TenantId:     opts.TenantId,
		Name:         opts.Name,
		Description:  opts.Description,
		AdminStateUp: opts.AdminStateUp,
		PolicyId:     opts.PolicyId,
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

// Get retrieves a particular firewall based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOpts contains the values used when updating a firewall.
type UpdateOpts struct {
	// Name of the firewall.
	Name         string
	Description  string
	AdminStateUp bool
	Status       string
	PolicyId     string
}

// Update allows firewalls to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type firewall struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		AdminStateUp bool   `json:"admin_state_up,omitempty"`
		PolicyId     string `json:"firewall_policy_id,omitempty"`
	}
	type request struct {
		Firewall firewall `json:"firewall"`
	}

	reqBody := request{Firewall: firewall{
		Name:         opts.Name,
		Description:  opts.Description,
		AdminStateUp: opts.AdminStateUp,
		PolicyId:     opts.PolicyId,
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

// Delete will permanently delete a particular firewall based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
