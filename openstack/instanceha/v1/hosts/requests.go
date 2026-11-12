package hosts

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add host listing parameters.
type ListOptsBuilder interface{ ToHostListQuery() (string, error) }

// ListOpts specifies filters and pagination options for listing hosts belonging to a failover segment.
type ListOpts struct {
	// Limit is the maximum number of resources returned per page.
	Limit int `q:"limit"`
	// Marker is the numeric database ID of the last resource in the previous page, encoded as a string.
	Marker string `q:"marker"`
	// SortKey is the resource field used to sort the results.
	SortKey string `q:"sort_key"`
	// SortDir is the sort direction: asc or desc.
	SortDir string `q:"sort_dir"`
	// Type filters by the host type, such as COMPUTE_HOST.
	Type string `q:"type"`
	// Reserved filters by whether the host is reserved for recovery.
	Reserved *bool `q:"reserved"`
	// OnMaintenance filters by whether the host is in maintenance mode.
	OnMaintenance *bool `q:"on_maintenance"`
}

// ToHostListQuery formats ListOpts as a URL query string.
func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a pager over hosts belonging to a failover segment. segmentID is the segment UUID.
func List(client *gophercloud.ServiceClient, segmentID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, segmentID)
	if opts != nil {
		query, err := opts.ToHostListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(result pagination.PageResult) pagination.Page {
		page := HostPage{MarkerPageBase: pagination.MarkerPageBase{PageResult: result}}
		page.Owner = &page
		return page
	})
}

// Get retrieves a host by UUID. segmentID is the segment UUID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, segmentID string, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, segmentID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
