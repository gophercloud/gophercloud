package vmoves

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add VM move listing parameters.
type ListOptsBuilder interface{ ToVMoveListQuery() (string, error) }

// ListOpts specifies filters and pagination options for listing VM moves belonging to a COMPUTE_HOST notification.
type ListOpts struct {
	// Limit is the maximum number of resources returned per page.
	Limit int `q:"limit"`
	// Marker is the numeric database ID of the last resource in the previous page, encoded as a string.
	Marker string `q:"marker"`
	// SortKey is the resource field used to sort the results.
	SortKey string `q:"sort_key"`
	// SortDir is the sort direction: asc or desc.
	SortDir string `q:"sort_dir"`
	// Status filters by the resource processing status.
	Status string `q:"status"`
	// Type filters by the VM move type, such as evacuation or migration.
	Type string `q:"type"`
}

// ToVMoveListQuery formats ListOpts as a URL query string.
func (opts ListOpts) ToVMoveListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a pager over VM moves belonging to a COMPUTE_HOST notification.
// notificationID is the notification UUID. Requires microversion 1.3 or later.
func List(client *gophercloud.ServiceClient, notificationID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, notificationID)
	if opts != nil {
		query, err := opts.ToVMoveListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(result pagination.PageResult) pagination.Page {
		page := VMovePage{MarkerPageBase: pagination.MarkerPageBase{PageResult: result}}
		page.Owner = &page
		return page
	})
}

// Get retrieves a VM move by UUID.
// notificationID is the parent notification UUID. Requires microversion 1.3 or later.
func Get(ctx context.Context, client *gophercloud.ServiceClient, notificationID string, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, notificationID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
