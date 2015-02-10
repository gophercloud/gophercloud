package instances

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/gophercloud/pagination"
)

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
	Databases os.DatabasesOpts

	// A slice of user information options.
	Users os.UsersOpts

	// ID of the configuration group to associate with the instance. Optional.
	ConfigID string

	// Options to configure the type of datastore the instance will use. This is
	// optional, and if excluded will default to MySQL.
	Datastore *DatastoreOpts

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
	instance, err := os.CreateOpts{
		FlavorRef: opts.FlavorRef,
		Size:      opts.Size,
		Name:      opts.Name,
		Databases: opts.Databases,
		Users:     opts.Users,
	}.ToInstanceCreateMap()

	if err != nil {
		return nil, err
	}

	instance = instance["instance"].(map[string]interface{})

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

	if opts.RestorePoint != "" {
		instance["restorePoint"] = opts.RestorePoint
	}

	return map[string]interface{}{"instance": instance}, nil
}

// Create will provision a new Database instance.
func Create(client *gophercloud.ServiceClient, opts os.CreateOptsBuilder) CreateResult {
	return CreateResult{os.Create(client, opts)}
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	return os.List(client)
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	return GetResult{os.Get(client, id)}
}
