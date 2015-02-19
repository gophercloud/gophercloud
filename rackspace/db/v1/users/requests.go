package users

import (
	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
)

/*
ChangePassword changes the password for one or more users. For example, to
change the respective passwords for two users:

	opts := os.BatchCreateOpts{
		os.CreateOpts{Name: "db_user_1", Password: "new_password_1"},
		os.CreateOpts{Name: "db_user_2", Password: "new_password_2"},
	}

	ChangePassword(client, "instance_id", opts)
*/
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

// Update will modify the attributes of a specified user. Attributes that can
// be updated are: user name, password, and host.
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

// Get will retrieve the details for a particular user.
func Get(client *gophercloud.ServiceClient, instanceID, userName string) GetResult {
	var res GetResult

	_, res.Err = client.Request("GET", userURL(client, instanceID, userName), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})

	return res
}

// ListAccess will list all of the databases a user has access to.
func ListAccess(client *gophercloud.ServiceClient, instanceID, userName string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return AccessPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, dbsURL(client, instanceID, userName), pageFn)
}

/*
GrantAccess for the specified user to one or more databases on a specified
instance. For example, to add a user to multiple databases:

	opts := db.BatchCreateOpts{
		db.CreateOpts{Name: "database_1"},
		db.CreateOpts{Name: "database_3"},
		db.CreateOpts{Name: "database_19"},
	}

	GrantAccess(client, "instance_id", "user_name", opts)
*/
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

/*
RevokeAccess will revoke access for the specified user to one or more databases
on a specified instance. For example, to remove a user's access to multiple
databases:

	opts := db.BatchCreateOpts{
		db.CreateOpts{Name: "database_1"},
		db.CreateOpts{Name: "database_3"},
		db.CreateOpts{Name: "database_19"},
	}

	RevokeAccess(client, "instance_id", "user_name", opts)
*/
func RevokeAccess(client *gophercloud.ServiceClient, instanceID, userName, dbName string) RevokeAccessResult {
	var res RevokeAccessResult

	_, res.Err = client.Request("DELETE", dbURL(client, instanceID, userName, dbName), gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return res
}
