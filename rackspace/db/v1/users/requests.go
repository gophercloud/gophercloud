package users

import (
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

	_, res.Err = client.Request("PUT", baseURL(client, instanceID), gophercloud.RequestOpts{
		JSONBody: &reqBody,
		OkCodes:  []int{202},
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

	_, res.Err = client.Request("PUT", userURL(client, instanceID, userName), gophercloud.RequestOpts{
		JSONBody: &reqBody,
		OkCodes:  []int{202},
	})

	return res
}

func Get(client *gophercloud.ServiceClient, instanceID, userName string) GetResult {
	var res GetResult

	_, res.Err = client.Request("GET", userURL(client, instanceID, userName), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
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

	_, res.Err = client.Request("PUT", dbsURL(client, instanceID, userName), gophercloud.RequestOpts{
		JSONBody: &reqBody,
		OkCodes:  []int{202},
	})

	return res
}

func RevokeAccess(client *gophercloud.ServiceClient, instanceID, userName, dbName string) RevokeAccessResult {
	var res RevokeAccessResult

	_, res.Err = client.Request("DELETE", dbURL(client, instanceID, userName, dbName), gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return res
}
