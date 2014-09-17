package extensions

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func Get(c *gophercloud.ServiceClient, name string) (*Extension, error) {
	var ext Extension
	_, err := perigee.Request("GET", ExtensionURL(c, name), perigee.Options{
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

func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, ListExtensionURL(c), func(r pagination.LastHTTPResponse) pagination.Page {
		return ExtensionPage{pagination.SinglePageBase(r)}
	})
}
