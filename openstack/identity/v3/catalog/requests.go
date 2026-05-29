package catalog

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List enumerates the services available to a specific user.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServiceCatalogPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
