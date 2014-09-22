package extensions

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Get retrieves information for a specific extension using its alias. If no
// extension exists with this alias, an error will be returned.
func Get(c *gophercloud.ServiceClient, alias string) (*Extension, error) {
	var ext Extension
	_, err := perigee.Request("GET", extensionURL(c, alias), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results: &struct {
			Extension *Extension `json:"extension"`
		}{&ext},
		OkCodes: []int{200},
	})

	if err != nil {
		return nil, err
	}
	return &ext, nil
}

// List returns a Pager which allows you to iterate over the full collection of
// extensions. It does not accept query parameters.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listExtensionURL(c), func(r pagination.LastHTTPResponse) pagination.Page {
		return ExtensionPage{pagination.SinglePageBase(r)}
	})
}
