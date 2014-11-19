package secgroups

import (
	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func commonList(client *gophercloud.ServiceClient, url string) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return SecurityGroupPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPage)
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	return commonList(client, rootURL(client))
}

func ListByServer(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	return commonList(client, listByServerURL(client, serverID))
}

type CreateOpts struct {
	// Optional - the name of your security group. If no value provided, null
	// will be set.
	Name string `json:"name,omitempty"`

	// Optional - the description of your security group. If no value provided,
	// null will be set.
	Description string `json:"description,omitempty"`
}

func Create(client *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var result CreateResult

	reqBody := struct {
		CreateOpts `json:"security_group"`
	}{opts}

	_, result.Err = perigee.Request("POST", rootURL(client), perigee.Options{
		Results:     &result.Body,
		ReqBody:     &reqBody,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return result
}
