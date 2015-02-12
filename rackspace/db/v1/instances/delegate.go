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

// Create asynchronously provisions a new database instance. It requires the
// user to specify a flavor and a volume size. The API service then provisions
// the instance with the requested flavor and sets up a volume of the specified
// size, which is the storage for the database instance.
//
// Although this call only allows the creation of 1 instance per request, you
// can create an instance with multiple databases and users. The default
// binding for a MySQL instance is port 3306.
func Create(client *gophercloud.ServiceClient, opts os.CreateOptsBuilder) CreateResult {
	return CreateResult{os.Create(client, opts)}
}

// List retrieves the status and information for all database instances.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return os.List(client)
}

// Get retrieves the status and information for a specified database instance.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	return GetResult{os.Get(client, id)}
}

// Delete permanently destroys the database instance.
func Delete(client *gophercloud.ServiceClient, id string) os.DeleteResult {
	return os.Delete(client, id)
}

// EnableRootUser enables the login from any host for the root user and
// provides the user with a generated root password.
func EnableRootUser(client *gophercloud.ServiceClient, id string) os.UserRootResult {
	return os.EnableRootUser(client, id)
}

// IsRootEnabled checks an instance to see if root access is enabled. It returns
// True if root user is enabled for the specified database instance or False
// otherwise.
func IsRootEnabled(client *gophercloud.ServiceClient, id string) (bool, error) {
	return os.IsRootEnabled(client, id)
}

// RestartService will restart only the MySQL Instance. Restarting MySQL will
// erase any dynamic configuration settings that you have made within MySQL.
// The MySQL service will be unavailable until the instance restarts.
func RestartService(client *gophercloud.ServiceClient, id string) os.ActionResult {
	return os.RestartService(client, id)
}

// ResizeInstance changes the memory size of the instance, assuming a valid
// flavorRef is provided. It will also restart the MySQL service.
func ResizeInstance(client *gophercloud.ServiceClient, id, flavorRef string) os.ActionResult {
	return os.ResizeInstance(client, id, flavorRef)
}

// ResizeVolume will resize the attached volume for an instance. It supports
// only increasing the volume size and does not support decreasing the size.
// The volume size is in gigabytes (GB) and must be an integer.
func ResizeVolume(client *gophercloud.ServiceClient, id string, size int) os.ActionResult {
	return os.ResizeVolume(client, id, size)
}
