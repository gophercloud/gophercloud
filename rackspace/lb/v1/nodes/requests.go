package nodes

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient, loadBalancerID int, limit *int) pagination.Pager {
	url := rootURL(client, loadBalancerID)

	if limit != nil {
		url += fmt.Sprintf("?limit=%d", limit)
	}

	createPageFn := func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, url, createPageFn)
}
