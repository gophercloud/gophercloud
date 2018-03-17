package policytypes

import (
	//"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	// Get retrieves details of a single policy type. Use ExtractPolicyType to convert its
	// result into a PolicyType.
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListDetail instructs OpenStack to provide a list of profile types.
func ListDetail(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PolicyTypePage{pagination.LinkedPageBase{PageResult: r}}
	})
}
