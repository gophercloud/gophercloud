package apiversions

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List lists all the Cinder API versions available to end-users.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return APIVersionPage{pagination.SinglePageBase(r)}
	})
}

// Get will retrieve the volume type with the provided ID. To extract the volume
// type from the result, call the Extract method on the GetResult.
func Get(ctx context.Context, client *gophercloud.ServiceClient, v string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, v), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
