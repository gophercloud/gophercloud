package extensions

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	common "github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Get retrieves information for a specific extension using its alias.
func Get(ctx context.Context, c *gophercloud.ServiceClient, alias string) common.GetResult {
	return common.Get(ctx, c, alias)
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return common.List(c).WithPageCreator(func(r pagination.PageResult) pagination.Page {
		return ExtensionPage{
			ExtensionPage: common.ExtensionPage{SinglePageBase: pagination.SinglePageBase(r)},
		}
	})
}
