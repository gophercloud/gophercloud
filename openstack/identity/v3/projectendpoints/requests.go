package projectendpoints

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type CreateOptsBuilder interface {
	ToEndpointCreateMap() (map[string]any, error)
}

// Create inserts a new Endpoint association to a project.
func Create(ctx context.Context, client *gophercloud.ServiceClient, projectID, endpointID string) (r CreateResult) {
	resp, err := client.Put(ctx, createURL(client, projectID, endpointID), nil, nil, &gophercloud.RequestOpts{OkCodes: []int{204}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List enumerates endpoints in a paginated collection, optionally filtered
// by ListOpts criteria.
func List(client *gophercloud.ServiceClient, projectID string) pagination.Pager {
	u := listURL(client, projectID)
	return pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return EndpointPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete removes an endpoint from the service catalog.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, projectID string, endpointID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, projectID, endpointID), &gophercloud.RequestOpts{OkCodes: []int{204}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
