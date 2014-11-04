package vips

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List is the operation responsible for returning a paginated collection of
// load balancer virtual IP addresses.
func List(client *gophercloud.ServiceClient, loadBalancerID int) pagination.Pager {
	url := rootURL(client, loadBalancerID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VIPPage{pagination.SinglePageBase(r)}
	})
}
