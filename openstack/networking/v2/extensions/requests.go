package extensions

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
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

func List(c *gophercloud.ServiceClient) gophercloud.Pager {
	return gophercloud.NewLinkedPager(c, ListExtensionURL(c))
}
