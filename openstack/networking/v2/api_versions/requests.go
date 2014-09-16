package networks

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

func List(c *gophercloud.ServiceClient) (*APIVersionsList, error) {
	var resp APIVersionsList
	_, err := perigee.Request("GET", APIVersionsURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func Get(c *gophercloud.ServiceClient, v string) (*APIInfoList, error) {
	var resp APIInfoList
	_, err := perigee.Request("GET", APIInfoURL(c, v), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
