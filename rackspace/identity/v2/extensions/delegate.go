package extensions

import (
	"github.com/rackspace/gophercloud"
	common "github.com/rackspace/gophercloud/openstack/common/extensions"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/extensions"
	"github.com/rackspace/gophercloud/pagination"
)

// ExtensionPage is the page returned by a pager when traversing over a collection of extensions.
type ExtensionPage struct {
	common.ExtensionPage
}

// IsEmpty checks whether an ExtensionPage struct is empty.
func (r ExtensionPage) IsEmpty() (bool, error) {
	is, err := ExtractExtensions(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractExtensions accepts a Page struct, specifically an ExtensionPage struct, and extracts the
// elements into a slice of os.Extension structs.
func ExtractExtensions(page pagination.Page) ([]os.Extension, error) {
	exts, err := common.ExtractExtensions(page.(ExtensionPage).ExtensionPage)
	if err != nil {
		return nil, err
	}
	results := make([]os.Extension, len(exts))
	for i, ext := range exts {
		results[i] = os.Extension{Extension: ext}
	}
	return results, nil
}

// Get retrieves information for a specific extension using its alias.
func Get(c *gophercloud.ServiceClient, alias string) os.GetResult {
	return os.Get(c, alias)
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	pager := os.List(c)
	pager.CreatePage = func(r pagination.LastHTTPResponse) pagination.Page {
		return ExtensionPage{
			ExtensionPage: common.ExtensionPage{SinglePageBase: pagination.SinglePageBase(r)},
		}
	}
	return pager
}
