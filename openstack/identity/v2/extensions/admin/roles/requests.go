package roles

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, rootURL(client), createPage)
}

func AddUserRole(client *gophercloud.ServiceClient, tenantID, userID, roleID string) UserRoleResult {
	var result UserRoleResult

	_, result.Err = perigee.Request("PUT", userRoleURL(client, tenantID, userID, roleID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{201},
	})

	return result
}

func DeleteUserRole(client *gophercloud.ServiceClient, tenantID, userID, roleID string) UserRoleResult {
	var result UserRoleResult

	_, result.Err = perigee.Request("DELETE", userRoleURL(client, tenantID, userID, roleID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})

	return result
}
