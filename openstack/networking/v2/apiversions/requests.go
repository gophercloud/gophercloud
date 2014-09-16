package apiversions

import "github.com/rackspace/gophercloud"

func ListVersions(c *gophercloud.ServiceClient) gophercloud.Pager {
	return gophercloud.NewLinkedPager(c, APIVersionsURL(c))
}

func ListVersionResources(c *gophercloud.ServiceClient, v string) gophercloud.Pager {
	return gophercloud.NewLinkedPager(c, APIInfoURL(c, v))
}
