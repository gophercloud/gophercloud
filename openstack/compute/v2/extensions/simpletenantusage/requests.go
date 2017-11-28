package simpletenantusage

import (
	"net/url"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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
		return SimpleSingleTenantUsagePage{pagination.SinglePageBase(r)}
	})
}

// GetOpts are options for fetching the tenant usage for a Tenant.
type GetOpts struct {
	// The ending time to calculate usage statistics on compute and storage resources.
	End *time.Time `json:"end,omitempty" q:"end"`

	// The beginning time to calculate usage statistics on compute and storage resources.
	Start *time.Time `json:"start,omitempty" q:"start"`
}

// GetOptsBuilder is an interface for GetOpts structs
type GetOptsBuilder interface {
	// Extra specific name to prevent collisions with interfaces for other usages
	ToSimpleTenantUsageGetMap() (string, error)
}

// ToSimpleTenantUsageGetMap converts the options into URL-encoded query string
// arguments.
func (opts GetOpts) ToSimpleTenantUsageGetMap() (string, error) {
	params := make(url.Values)
	if opts.Start != nil {
		params.Add("start", opts.Start.Format(gophercloud.RFC3339MilliNoZ))
	}

	if opts.End != nil {
		params.Add("end", opts.End.Format(gophercloud.RFC3339MilliNoZ))
	}

	q := &url.URL{RawQuery: params.Encode()}
	return q.String(), nil
}
