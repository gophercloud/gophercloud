package extensions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get retrieves information for a specific extension using its alias.
func Get(c *gophercloud.ServiceClient, alias string) GetResult {
	var res GetResult
	_, res.Err = c.Get(ExtensionURL(c, alias), &res.Body, nil)
	return res
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, ListExtensionURL(c), func(r pagination.PageResult) pagination.Page {
		return ExtensionPage{pagination.SinglePageBase(r)}
	})
}
