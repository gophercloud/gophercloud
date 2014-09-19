package apiversions

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListVersions lists all the Neutron API versions available to end-users
func ListVersions(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, APIVersionsURL(c), func(r pagination.LastHTTPResponse) pagination.Page {
		return APIVersionPage{pagination.SinglePageBase(r)}
	})
}

// ListVersionResources lists all of the different API resources for a particular
// API versions. Typical resources for Neutron might be: networks, subnets, etc.
func ListVersionResources(c *gophercloud.ServiceClient, v string) pagination.Pager {
	return pagination.NewPager(c, APIInfoURL(c, v), func(r pagination.LastHTTPResponse) pagination.Page {
		return APIVersionResourcePage{pagination.SinglePageBase(r)}
	})
}
