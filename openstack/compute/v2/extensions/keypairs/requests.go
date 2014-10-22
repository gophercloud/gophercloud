package keypairs

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of KeyPairs.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return KeyPairPage{pagination.SinglePageBase(r)}
	})
}
