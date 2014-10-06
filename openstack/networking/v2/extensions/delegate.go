package extensions

import (
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

// ExtractExtensions interprets a Page as a slice of Extensions.
func ExtractExtensions(page pagination.Page) ([]Extension, error) {
	inner, err := common.ExtractExtensions(page)
	if err != nil {
		return nil, err
	}
	outer := make([]Extension, len(inner))
	for index, ext := range inner {
		outer[index] = Extension{ext}
	}
	return outer, nil
}

// rebased is a temporary workaround to isolate changes to this package. FIXME: set ResourceBase
// in the NewNetworkV2 method and remove the version string from URL generation methods in
// networking resources.
func rebased(c *gophercloud.ServiceClient) *gophercloud.ServiceClient {
	var r = *c
	r.ResourceBase = c.Endpoint + "v2.0/"
	return &r
}

// Get retrieves information for a specific extension using its alias.
func Get(c *gophercloud.ServiceClient, alias string) GetResult {
	return GetResult{common.Get(rebased(c), alias)}
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return common.List(rebased(c))
}
