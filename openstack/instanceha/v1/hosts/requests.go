package hosts

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type ListOptsBuilder interface{ ToHostListQuery() (string, error) }

type ListOpts struct {
	Limit             int    `q:"limit"`
	Marker            string `q:"marker"`
	SortKey           string `q:"sort_key"`
	SortDir           string `q:"sort_dir"`
	Name              string `q:"name"`
	Type              string `q:"type"`
	ControlAttributes string `q:"control_attributes"`
	Reserved          *bool  `q:"reserved"`
	OnMaintenance     *bool  `q:"on_maintenance"`
}

func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

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

func Get(ctx context.Context, client *gophercloud.ServiceClient, segmentID string, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, segmentID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
