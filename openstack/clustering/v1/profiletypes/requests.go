package profiletypes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	// Get retrieves details of a single profile type. Use ExtractProfileType to convert its
	// result into a ProfileType.
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListDetail instructs OpenStack to provide a list of profile types.
func ListDetail(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProfileTypePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListOperationDetail instructs OpenStack to provide a list of profile type operations.
func ListOperationDetail(client *gophercloud.ServiceClient, id string) pagination.Pager {
	url := listOperationURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return OperationPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
