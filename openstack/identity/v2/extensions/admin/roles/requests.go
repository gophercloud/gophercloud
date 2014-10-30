package roles

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, rootURL(client), createPage)
}

func AddRoleToUser(client *gophercloud.ServiceClient, tenantID, userID, roleID string) AddRoleResult {
	var result AddRoleResult

	_, result.Err = perigee.Request("PUT", userRoleURL(client, tenantID, userID, roleID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{201},
	})

	return result
}
