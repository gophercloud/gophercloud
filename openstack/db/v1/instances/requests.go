package instances

import (
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
)

// CreateOptsBuilder is the top-level interface for create options.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// UserOpts is the struct responsible for configuring a user; often in the
// context of an instance.
type UserOpts struct {
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

func (opts UserOpts) ToMap() (map[string]interface{}, error) {
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

type UsersOpts []UserOpts

func (opts UsersOpts) ToMap() ([]map[string]interface{}, error) {
	var users []map[string]interface{}
	for _, opt := range opts {
		user, err := opt.ToMap()
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// CreateOpts is the struct responsible for configuring a new database instance.
type CreateOpts struct {
	// Either the integer UUID (in string form) of the flavor, or its URI
	// reference as specified in the response from the List() call. Required.
	FlavorRef string

	// Specifies the volume size in gigabytes (GB). The value must be between 1
	// and 300. Required.
	Size int

	// Name of the instance to create. The length of the name is limited to
	// 255 characters and any characters are permitted. Optional.
	Name string

	// A slice of database information options.
	Databases db.BatchCreateOpts

	// A slice of user information options.
	Users UsersOpts
}

func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	if opts.Size > 300 || opts.Size < 1 {
		return nil, fmt.Errorf("Size (GB) must be between 1-300")
	}
	if opts.FlavorRef == "" {
		return nil, fmt.Errorf("FlavorRef is a required field")
	}

	instance := map[string]interface{}{
		"volume":    map[string]int{"size": opts.Size},
		"flavorRef": opts.FlavorRef,
	}

	if opts.Name != "" {
		instance["name"] = opts.Name
	}
	if len(opts.Databases) > 0 {
		dbs, err := opts.Databases.ToDBCreateMap()
		if err != nil {
			return nil, err
		}
		instance["databases"] = dbs["databases"]
	}
	if len(opts.Users) > 0 {
		users, err := opts.Users.ToMap()
		if err != nil {
			return nil, err
		}
		instance["users"] = users
	}

	return map[string]interface{}{"instance": instance}, nil
}

// Create asynchronously provisions a new database instance. It requires the
// user to specify a flavor and a volume size. The API service then provisions
// the instance with the requested flavor and sets up a volume of the specified
// size, which is the storage for the database instance.
//
// Although this call only allows the creation of 1 instance per request, you
// can create an instance with multiple databases and users. The default
// binding for a MySQL instance is port 3306.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToInstanceCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	resp, err := perigee.Request("POST", baseURL(client), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// List retrieves the status and information for all database instances.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPageFn := func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, baseURL(client), createPageFn)
}

// Get retrieves the status and information for a specified database instance.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult

	resp, err := perigee.Request("GET", resourceURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// Delete permanently destroys the database instance.
func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult

	resp, err := perigee.Request("DELETE", resourceURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// EnableRootUser enables the login from any host for the root user and
// provides the user with a generated root password.
func EnableRootUser(client *gophercloud.ServiceClient, id string) UserRootResult {
	var res UserRootResult

	resp, err := perigee.Request("POST", userRootURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// IsRootEnabled checks an instance to see if root access is enabled. It returns
// True if root user is enabled for the specified database instance or False
// otherwise.
func IsRootEnabled(client *gophercloud.ServiceClient, id string) (bool, error) {
	var res gophercloud.Result

	_, err := perigee.Request("GET", userRootURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res.Body.(map[string]interface{})["rootEnabled"] == true, err
}

// RestartService will restart only the MySQL Instance. Restarting MySQL will
// erase any dynamic configuration settings that you have made within MySQL.
// The MySQL service will be unavailable until the instance restarts.
func RestartService(client *gophercloud.ServiceClient, id string) ActionResult {
	var res ActionResult

	resp, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     map[string]bool{"restart": true},
		OkCodes:     []int{202},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// ResizeInstance changes the memory size of the instance, assuming a valid
// flavorRef is provided. It will also restart the MySQL service.
func ResizeInstance(client *gophercloud.ServiceClient, id, flavorRef string) ActionResult {
	var res ActionResult

	reqBody := map[string]map[string]string{
		"resize": map[string]string{
			"flavorRef": flavorRef,
		},
	}

	resp, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

// ResizeVolume will resize the attached volume for an instance. It supports
// only increasing the volume size and does not support decreasing the size.
// The volume size is in gigabytes (GB) and must be an integer.
func ResizeVolume(client *gophercloud.ServiceClient, id string, size int) ActionResult {
	var res ActionResult

	reqBody := map[string]map[string]map[string]int{
		"resize": map[string]map[string]int{
			"volume": map[string]int{"size": size},
		},
	}

	resp, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}
