package extensions

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Get retrieves information for a specific extension using its alias.
func Get(c *gophercloud.ServiceClient, alias string) GetResult {
	var res GetResult
	_, err := perigee.Request("GET", extensionURL(c, alias), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	res.Err = err
	return res
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listExtensionURL(c), func(r pagination.LastHTTPResponse) pagination.Page {
		return ExtensionPage{pagination.SinglePageBase(r)}
	})
}
