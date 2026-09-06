package hosts

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List retrieves a list of compute hosts.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return HostPage{pagination.SinglePageBase(r)}
	})
}
