package routers

import (
	"strconv"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

const (
	version      = "v2.0"
	resourcePath = "routers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(version, resourcePath, id)
}

type ListOpts struct {
	ID           string
	Name         string
	AdminStateUp *bool
	Status       string
	TenantID     string
	Limit        int
	Marker       string
	SortKey      string
	SortDir      string
}

func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q := make(map[string]string)
	if opts.ID != "" {
		q["id"] = opts.ID
	}
	if opts.Name != "" {
		q["name"] = opts.Name
	}
	if opts.AdminStateUp != nil {
		q["admin_state_up"] = strconv.FormatBool(*opts.AdminStateUp)
	}
	if opts.Status != "" {
		q["status"] = opts.Status
	}
	if opts.TenantID != "" {
		q["tenant_id"] = opts.TenantID
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

	u := rootURL(c) + utils.BuildQuery(q)
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return RouterPage{pagination.LinkedPageBase(r)}
	})
}
