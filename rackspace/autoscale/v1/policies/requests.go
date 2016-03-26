package policies

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns all scaling policies for a group.
func List(client *gophercloud.ServiceClient, groupID string) pagination.Pager {
	url := listURL(client, groupID)

	createPageFn := func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPageFn)
}
