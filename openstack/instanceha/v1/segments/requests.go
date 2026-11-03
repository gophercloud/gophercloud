package segments

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add failover segment listing parameters.
type ListOptsBuilder interface{ ToSegmentListQuery() (string, error) }

// ListOpts specifies filters and pagination options for listing failover segments.
type ListOpts struct {
	// Limit is the maximum number of resources returned per page.
	Limit int `q:"limit"`
	// Marker is the numeric database ID of the last resource in the previous page, encoded as a string.
	Marker string `q:"marker"`
	// SortKey is the resource field used to sort the results.
	SortKey string `q:"sort_key"`
	// SortDir is the sort direction: asc or desc.
	SortDir string `q:"sort_dir"`
	// RecoveryMethod filters by the recovery strategy: auto, reserved_host, auto_priority, or rh_priority.
	RecoveryMethod string `q:"recovery_method"`
	// ServiceType filters by the service managed by the segment, such as COMPUTE.
	ServiceType string `q:"service_type"`
	// Enabled filters by whether recovery is enabled for the segment. Requires microversion 1.2 or later.
	Enabled *bool `q:"enabled"`
}

// ToSegmentListQuery formats ListOpts as a URL query string.
func (opts ListOpts) ToSegmentListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a pager over failover segments.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToSegmentListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(result pagination.PageResult) pagination.Page {
		page := SegmentPage{MarkerPageBase: pagination.MarkerPageBase{PageResult: result}}
		page.Owner = &page
		return page
	})
}

// Get retrieves a segment by UUID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
