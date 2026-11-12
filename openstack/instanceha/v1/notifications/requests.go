package notifications

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add notification listing parameters.
type ListOptsBuilder interface{ ToNotificationListQuery() (string, error) }

// ListOpts specifies filters and pagination options for listing failure notifications.
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
	// Type filters by the event type, such as COMPUTE_HOST, PROCESS, or VM.
	Type string `q:"type"`
	// SourceHostUUID filters by the UUID of the host that generated the notification.
	SourceHostUUID string `q:"source_host_uuid"`
	// GeneratedSince filters notifications generated at or after this ISO 8601 timestamp.
	GeneratedSince string `q:"generated-since"`
}

// ToNotificationListQuery formats ListOpts as a URL query string.
func (opts ListOpts) ToNotificationListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a pager over failure notifications.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToNotificationListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(result pagination.PageResult) pagination.Page {
		page := NotificationPage{MarkerPageBase: pagination.MarkerPageBase{PageResult: result}}
		page.Owner = &page
		return page
	})
}

// Get retrieves a notification by UUID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
