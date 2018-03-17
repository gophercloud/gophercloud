package apiversion

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListDetail instructs OpenStack to provide a list of cluster.
func ListDetail(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VersionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
