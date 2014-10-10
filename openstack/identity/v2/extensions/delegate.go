package extensions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	common "github.com/rackspace/gophercloud/openstack/common/extensions"
	"github.com/rackspace/gophercloud/pagination"
)

// Extension is a single OpenStack extension.
type Extension struct {
	common.Extension
}

// GetResult wraps a GetResult from common.
type GetResult struct {
	common.GetResult
}

// ExtensionPage is a single page of Extension results.
type ExtensionPage struct {
	common.ExtensionPage
}

// IsEmpty returns true if the current page contains at least one Extension.
func (page ExtensionPage) IsEmpty() (bool, error) {
	is, err := ExtractExtensions(page)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractExtensions accepts a Page struct, specifically an ExtensionPage struct, and extracts the
// elements into a slice of Extension structs.
func ExtractExtensions(page pagination.Page) ([]Extension, error) {
	// Identity v2 adds an intermediate "values" object.

	type extension struct {
		Updated     string        `mapstructure:"updated"`
		Name        string        `mapstructure:"name"`
		Namespace   string        `mapstructure:"namespace"`
		Alias       string        `mapstructure:"alias"`
		Description string        `mapstructure:"description"`
		Links       []interface{} `mapstructure:"links"`
	}

	var resp struct {
		Extensions struct {
			Values []extension `mapstructure:"values"`
		} `mapstructure:"extensions"`
	}

	err := mapstructure.Decode(page.(ExtensionPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	exts := make([]Extension, len(resp.Extensions.Values))
	for i, original := range resp.Extensions.Values {
		exts[i] = Extension{common.Extension{
			Updated:     original.Updated,
			Name:        original.Name,
			Namespace:   original.Namespace,
			Alias:       original.Alias,
			Description: original.Description,
			Links:       original.Links,
		}}
	}

	return exts, err
}

// Get retrieves information for a specific extension using its alias.
func Get(c *gophercloud.ServiceClient, alias string) GetResult {
	return GetResult{common.Get(c, alias)}
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return common.List(c).WithPageCreator(func(r pagination.LastHTTPResponse) pagination.Page {
		return ExtensionPage{
			ExtensionPage: common.ExtensionPage{SinglePageBase: pagination.SinglePageBase(r)},
		}
	})
}
