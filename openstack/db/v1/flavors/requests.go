package flavors

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, listURL(client), createPage)
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var gr GetResult
	gr.Err = perigee.Get(getURL(client, id), perigee.Options{
		Results:     &gr.Body,
		MoreHeaders: client.AuthenticatedHeaders(),
	})
	return gr
}
