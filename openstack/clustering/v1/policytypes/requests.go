package policytypes

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List makes a request against the API to list policy types.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := policyTypeListURL(client)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PolicyTypePage{pagination.SinglePageBase(r)}
	})
}

// Get makes a request against the API to get details for a policy type.
func Get(ctx context.Context, client *gophercloud.ServiceClient, policyTypeName string) (r GetResult) {
	url := policyTypeGetURL(client, policyTypeName)

	resp, err := client.GetWithContext(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
