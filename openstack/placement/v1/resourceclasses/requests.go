package resourceclasses

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List retrieves a list of resource classes.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourceClassesPage{pagination.SinglePageBase(r)}
	})
}

// Get retrieves the resource class with the provided name.
func Get(ctx context.Context, client *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, name), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
