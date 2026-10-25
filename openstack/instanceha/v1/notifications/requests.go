package notifications

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type ListOptsBuilder interface{ ToNotificationListQuery() (string, error) }

type ListOpts struct {
	Limit          int    `q:"limit"`
	Marker         string `q:"marker"`
	SortKey        string `q:"sort_key"`
	SortDir        string `q:"sort_dir"`
	Status         string `q:"status"`
	Type           string `q:"type"`
	SourceHostUUID string `q:"source_host_uuid"`
	GeneratedSince string `q:"generated-since"`
}

func (opts ListOpts) ToNotificationListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

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

func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
