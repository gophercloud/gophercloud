package federation

import (
	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/pagination"
)

// ListMappings enumerates the mappings.
func ListMappings(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, mappingsRootURL(client), func(r pagination.PageResult) pagination.Page {
		return MappingsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
