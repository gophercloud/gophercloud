package apiversions

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func ListVersions(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, APIVersionsURL(c), func(r pagination.LastHTTPResponse) pagination.Page {
		return APIVersionPage{pagination.SinglePageBase(r)}
	})
}

func ListVersionResources(c *gophercloud.ServiceClient, v string) pagination.Pager {
	return pagination.NewPager(c, APIInfoURL(c, v), func(r pagination.LastHTTPResponse) pagination.Page {
		return APIVersionResourcePage{pagination.SinglePageBase(r)}
	})
}
