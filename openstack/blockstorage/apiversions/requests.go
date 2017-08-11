package apiversions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List lists all the Cinder API versions available to end-users.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return APIVersionPage{pagination.SinglePageBase(r)}
	})
}

// Get will retrieve a specific API version. To extract the version,
// call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, v string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, v), &r.Body, nil)
	return
}
