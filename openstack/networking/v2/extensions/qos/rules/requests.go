package rules

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToBandwidthLimitRuleListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the BandwidthLimitRules attributes you want to see returned.
// SortKey allows you to sort by a particular BandwidthLimitRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type BandwidthLimitRuleListOpts struct {
	ID           string `q:"id"`
	TenantID     string `q:"tenant_id"`
	MaxKBps      int    `q:"max_kbps"`
	MaxBurstKBps int    `q:"max_burst_kbps"`
	Direction    string `q:"direction"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
	Tags         string `q:"tags"`
	TagsAny      string `q:"tags-any"`
	NotTags      string `q:"not-tags"`
	NotTagsAny   string `q:"not-tags-any"`
}

// ToBandwidthLimitRuleListQuery formats a ListOpts into a query string.
func (opts BandwidthLimitRuleListOpts) ToBandwidthLimitRuleListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// BandwidthLimitRuleList returns a Pager which allows you to iterate over a collection of
// BandwidthLimitRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func BandwidthLimitRuleList(c *gophercloud.ServiceClient, policyID string, opts ListOptsBuilder) pagination.Pager {
	url := listBandwidthLimitRulesURL(c, policyID)
	fmt.Printf("URL: %s\n", url)
	if opts != nil {
		query, err := opts.ToBandwidthLimitRuleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BandwidthLimitRulePage{pagination.LinkedPageBase{PageResult: r}}

	})
}
