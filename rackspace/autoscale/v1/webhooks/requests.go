package webhooks

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns all webhooks for a scaling policy.
func List(client *gophercloud.ServiceClient, groupID, policyID string) pagination.Pager {
	url := listURL(client, groupID, policyID)

	createPageFn := func(r pagination.PageResult) pagination.Page {
		return WebhookPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPageFn)
}
