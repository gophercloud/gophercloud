package users

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
)

func ChangePassword(client *gophercloud.ServiceClient, instanceID string, opts os.BatchCreateOpts) UpdatePasswordsResult {
	var res UpdatePasswordsResult

	reqBody, err := opts.ToUserCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("PUT", baseURL(client, instanceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})

	return res
}

func Update(client *gophercloud.ServiceClient, instanceID, userName string, opts os.CreateOpts) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToMap()
	if err != nil {
		res.Err = err
		return res
	}
	reqBody = map[string]interface{}{"user": reqBody}

	_, res.Err = perigee.Request("PUT", userURL(client, instanceID, userName), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})

	return res
}

func Get(client *gophercloud.ServiceClient, instanceID, userName string) GetResult {
	var res GetResult

	_, res.Err = perigee.Request("GET", userURL(client, instanceID, userName), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

func ListAccess(client *gophercloud.ServiceClient, instanceID, userName string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return AccessPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, dbsURL(client, instanceID, userName), pageFn)
}

func GrantAccess(client *gophercloud.ServiceClient, instanceID, userName string, opts db.BatchCreateOpts) GrantAccessResult {
	var res GrantAccessResult

	reqBody, err := opts.ToDBCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("PUT", dbsURL(client, instanceID, userName), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})

	return res
}

func RevokeAccess(client *gophercloud.ServiceClient, instanceID, userName, dbName string) RevokeAccessResult {
	var res RevokeAccessResult

	_, res.Err = perigee.Request("DELETE", dbURL(client, instanceID, userName, dbName), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return res
}
