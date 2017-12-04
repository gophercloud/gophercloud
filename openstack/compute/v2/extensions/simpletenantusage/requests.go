package simpletenantusage

import (
	"net/url"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// GetSingleTenant returns simple tenant usage data about a specific tenant
func GetSingleTenant(client *gophercloud.ServiceClient, tenantID string, opts GetSingleTenantOptsBuilder) pagination.Pager {
	url := getTenantURL(client, tenantID)
	if opts != nil {
		query, err := opts.ToSingleTenantUsageGetMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SimpleSingleTenantUsagePage{pagination.SinglePageBase(r)}
	})
}

// GetSingleOpts are options for fetching the tenant usage for a Tenant.
type GetSingleTenantOpts struct {
	// The ending time to calculate usage statistics on compute and storage resources.
	End *time.Time `json:"end,omitempty" q:"end"`

	// The beginning time to calculate usage statistics on compute and storage resources.
	Start *time.Time `json:"start,omitempty" q:"start"`
}

// GetSingleTenantOptsBuilder is an interface for GetOpts structs
type GetSingleTenantOptsBuilder interface {
	// Extra specific name to prevent collisions with interfaces for other usages
	ToSingleTenantUsageGetMap() (string, error)
}

// ToSimpleTenantUsageGetMap converts the options into URL-encoded query string
// arguments.
func (opts GetSingleTenantOpts) ToSingleTenantUsageGetMap() (string, error) {
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
