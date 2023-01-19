package projectendpoints

import (
	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToEndpointCreateMap() (map[string]interface{}, error)
}

// Create inserts a new Endpoint association to a project.
func Create(client *gophercloud.ServiceClient, projectID, endpointID string) (r CreateResult) {
	resp, err := client.Put(createURL(client, projectID, endpointID), nil, nil, &gophercloud.RequestOpts{OkCodes: []int{204}})
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
func Delete(client *gophercloud.ServiceClient, projectID string, endpointID string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, projectID, endpointID), &gophercloud.RequestOpts{OkCodes: []int{204}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
