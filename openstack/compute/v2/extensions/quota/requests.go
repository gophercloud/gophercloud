package quotas

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return QuotaPage{pagination.SinglePageBase(r)}
	})
}

// Get returns public data about a previously created Quota.
func Get(client *gophercloud.ServiceClient, name string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getURL(client, name), &res.Body, nil)
	return res
}
