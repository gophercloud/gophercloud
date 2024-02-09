package profiletypes

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.GetWithContext(ctx, getURL(client, id), &r.Body,
		&gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List makes a request against the API to list profile types.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProfileTypePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListOps(client *gophercloud.ServiceClient, id string) pagination.Pager {
	url := listOpsURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return OperationPage{pagination.SinglePageBase(r)}
	})
}
