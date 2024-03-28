package roles

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List is the operation responsible for listing all available global roles
// that a user can adopt.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.SinglePageBase(r)}
	})
}

// AddUser is the operation responsible for assigning a particular role to
// a user. This is confined to the scope of the user's tenant - so the tenant
// ID is a required argument.
func AddUser(ctx context.Context, client *gophercloud.ServiceClient, tenantID, userID, roleID string) (r UserRoleResult) {
	resp, err := client.Put(ctx, userTenantRoleURL(client, tenantID, userID, roleID), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteUser is the operation responsible for deleting a particular role
// from a user. This is confined to the scope of the user's tenant - so the
// tenant ID is a required argument.
func DeleteUser(ctx context.Context, client *gophercloud.ServiceClient, tenantID, userID, roleID string) (r UserRoleResult) {
	resp, err := client.Delete(ctx, userTenantRoleURL(client, tenantID, userID, roleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
