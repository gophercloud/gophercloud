package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List will list all available hardware flavors that an instance can use. The
// operation is identical to the one supported by the Nova API, but without the
// "disk" property.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get will retrieve information for a specified hardware flavor.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
