package instances

import (
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// CreateOptsBuilder is the top-level interface for create options.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// DatastoreOpts represents the configuration for how an instance stores data.
type DatastoreOpts struct {
	Version string
	Type    string
}

func (opts DatastoreOpts) ToMap() (map[string]string, error) {
	return map[string]string{
		"version": opts.Version,
		"type":    opts.Type,
	}, nil
}

// DatabaseOpts is the struct responsible for configuring a database; often in
// the context of an instance.
type DatabaseOpts struct {
	// Specifies the name of the database. Optional.
	Name string

	// Set of symbols and encodings. Optional; the default character set is utf8.
	CharSet string

	// Set of rules for comparing characters in a character set. Optional; the
	// default value for collate is utf8_general_ci.
	Collate string
}

func (opts DatabaseOpts) ToMap() (map[string]string, error) {
	db := map[string]string{}
	if opts.Name != "" {
		db["name"] = opts.Name
	}
	if opts.CharSet != "" {
		db["character_set"] = opts.CharSet
	}
	if opts.Collate != "" {
		db["collate"] = opts.Collate
	}
	return db, nil
}

type DatabasesOpts []DatabaseOpts

func (opts DatabasesOpts) ToMap() ([]map[string]string, error) {
	var dbs []map[string]string
	for _, db := range opts {
		dbMap, err := db.ToMap()
		if err != nil {
			return dbs, err
		}
		dbs = append(dbs, dbMap)
	}
	return dbs, nil
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
	Databases []DatabaseOpts

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

	// ID of the configuration group to associate with the instance. Optional.
	ConfigID string

	// Options to configure the type of datastore the instance will use. This is
	// optional, and if excluded will default to MySQL.
	Datastore *DatastoreOpts

	// A slice of database information options.
	Databases DatabasesOpts

	// A slice of user information options.
	Users UsersOpts

	// Specifies the backup ID from which to restore the database instance. There
	// are some things to be aware of before using this field.  When you execute
	// the Restore Backup operation, a new database instance is created to store
	// the backup whose ID is specified by the restorePoint attribute. This will
	// mean that:
	// - All users, passwords and access that were on the instance at the time of
	// the backup will be restored along with the databases.
	// - You can create new users or databases if you want, but they cannot be
	// the same as the ones from the instance that was backed up.
	RestorePoint string
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
	if opts.ConfigID != "" {
		instance["configuration"] = opts.ConfigID
	}
	if opts.Datastore != nil {
		ds, err := opts.Datastore.ToMap()
		if err != nil {
			return nil, err
		}
		instance["datastore"] = ds
	}
	if len(opts.Databases) > 0 {
		dbs, err := opts.Databases.ToMap()
		if err != nil {
			return nil, err
		}
		instance["databases"] = dbs
	}
	if len(opts.Users) > 0 {
		users, err := opts.Users.ToMap()
		if err != nil {
			return nil, err
		}
		instance["users"] = users
	}
	if opts.RestorePoint != "" {
		instance["restorePoint"] = opts.RestorePoint
	}

	return map[string]interface{}{"instance": instance}, nil
}

// Create will provision a new Database instance.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToInstanceCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	resp, err := perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{200},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}
