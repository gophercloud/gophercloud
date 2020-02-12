package instanceactions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the API to list the servers actions.
func List(client *gophercloud.ServiceClient, id string) pagination.Pager {
	return pagination.NewPager(client, ListURL(client, id), func(r pagination.PageResult) pagination.Page {
		return InstanceActionPage{pagination.SinglePageBase(r)}
	})
}

// Get makes a request against the API to get a server action.
func Get(client *gophercloud.ServiceClient, serverID, requestID string) (r InstanceActionResult) {
	_, r.Err = client.Get(instanceActionsURL(client, serverID, requestID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
