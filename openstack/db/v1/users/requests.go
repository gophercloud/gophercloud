package users

import (
	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the struct responsible for configuring a new user; often in the
// context of an instance.
type CreateOpts struct {
	// Specifies a name for the user.
	Name string

	// Specifies a password for the user.
	Password string

	// An array of databases that this user will connect to. The `name` field is
	// the only requirement for each option.
	Databases db.BatchCreateOpts

	// Specifies the host from which a user is allowed to connect to the database.
	// Possible values are a string containing an IPv4 address or "%" to allow
	// connecting from any host. Optional; the default is "%".
	Host string
}

func (opts CreateOpts) ToMap() (map[string]interface{}, error) {
	user := map[string]interface{}{}

	if opts.Name != "" {
		user["name"] = opts.Name
	}
	if opts.Password != "" {
		user["password"] = opts.Password
	}
	if opts.Host != "" {
		user["host"] = opts.Host
	}

	var dbs []map[string]string
	for _, db := range opts.Databases {
		dbs = append(dbs, map[string]string{"name": db.Name})
	}
	if len(dbs) > 0 {
		user["databases"] = dbs
	}

	return user, nil
}

type BatchCreateOpts []CreateOpts

func (opts BatchCreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	var users []map[string]interface{}
	for _, opt := range opts {
		user, err := opt.ToMap()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return map[string]interface{}{"users": users}, nil
}

func Create(client *gophercloud.ServiceClient, instanceID string, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToUserCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("POST", baseURL(client, instanceID), gophercloud.RequestOpts{
		JSONBody: &reqBody,
		OkCodes:  []int{202},
	})

	return res
}

func List(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	createPageFn := func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, baseURL(client, instanceID), createPageFn)
}

func Delete(client *gophercloud.ServiceClient, instanceID, userName string) DeleteResult {
	var res DeleteResult

	_, res.Err = client.Request("DELETE", userURL(client, instanceID, userName), gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return res
}
