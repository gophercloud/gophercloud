package simpletenantusage

import (
	"time"

	"code.comcast.com/onecloud/gophercloud"
	"code.comcast.com/onecloud/gophercloud/pagination"
)

// Get returns simple tenant usage data about all tenants
func Get(client *gophercloud.ServiceClient, opts GetOptsBuilder) pagination.Pager {
	url := getURL(client)
	if opts != nil {
		query, err := opts.ToSimpleTenantUsageGetMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SimpleTenantUsagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetTenant returns simple tenant usage data about a specific tenant
func GetTenant(client *gophercloud.ServiceClient, tenantID string, opts GetOptsBuilder) pagination.Pager {
	url := getTenantURL(client, tenantID)
	if opts != nil {
		query, err := opts.ToSimpleTenantUsageGetMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SimpleTenantUsagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Options for fetching the tenant usage for a Tenant.
// All int-values are pointers so they can be nil if they are not needed
// you can use gopercloud.IntToPointer() for convenience
type GetOpts struct {
	// Set to 1 to get detailed (i.e. "server usages") tenant usage; 0 otherwise.
	Detailed *bool `json:"detailed,omitempty"`

	// The ending time to calculate usage statistics on compute and storage resources.
	End time.Time `json:"end,omitempty"`

	// The beginning time to calculate usage statistics on compute and storage resources.
	Start time.Time `json:"start,omitempty"`

	// Requests a page size of items. Calculate usage for the limited number of instances.
	// Use the limit parameter to make an initial limited request and use the last-seen instance
	// UUID from the response as the marker parameter value in a subsequent limited request.
	Limit *int `json:"limit,omitempty"`

	// The last-seen item based on an earlier Limit query.
	Marker string `json:"marker,omitempty"`
}

// Interface for GetOpts structs
type GetOptsBuilder interface {
	// Extra specific name to prevent collisions with interfaces for other usages
	ToSimpleTenantUsageGetMap() (string, error)
}

// Convert the options into URL-encoded query string arguments.
func (opts GetOpts) ToSimpleTenantUsageGetMap() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}
