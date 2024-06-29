package instances

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	db "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder is the top-level interface for create options.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]any, error)
}

// DatastoreOpts represents the configuration for how an instance stores data.
type DatastoreOpts struct {
	Version string `json:"version"`
	Type    string `json:"type"`
}

// ToMap converts a DatastoreOpts to a map[string]string (for a request body)
func (opts DatastoreOpts) ToMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// NetworkOpts is used within CreateOpts to control a new server's network attachments.
type NetworkOpts struct {
	// UUID of a nova-network to attach to the newly provisioned server.
	// Required unless Port is provided.
	UUID string `json:"net-id,omitempty"`

	// Port of a neutron network to attach to the newly provisioned server.
	// Required unless UUID is provided.
	Port string `json:"port-id,omitempty"`

	// V4FixedIP [optional] specifies a fixed IPv4 address to be used on this network.
	V4FixedIP string `json:"v4-fixed-ip,omitempty"`

	// V6FixedIP [optional] specifies a fixed IPv6 address to be used on this network.
	V6FixedIP string `json:"v6-fixed-ip,omitempty"`
}

// ToMap converts a NetworkOpts to a map[string]string (for a request body)
func (opts NetworkOpts) ToMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// CreateOpts is the struct responsible for configuring a new database instance.
type CreateOpts struct {
	// The availability zone of the instance.
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// ID of the configuration group that you want to attach to the instance.
	Configuration string `json:"configuration,omitempty"`
	// Either the integer UUID (in string form) of the flavor, or its URI
	// reference as specified in the response from the List() call. Required.
	FlavorRef string
	// Specifies the volume size in gigabytes (GB). The value must be between 1
	// and 300. Required.
	Size int
	// Specifies the volume type.
	VolumeType string
	// Name of the instance to create. The length of the name is limited to
	// 255 characters and any characters are permitted. Optional.
	Name string
	// A slice of database information options.
	Databases db.CreateOptsBuilder
	// A slice of user information options.
	Users users.CreateOptsBuilder
	// Options to configure the type of datastore the instance will use. This is
	// optional, and if excluded will default to MySQL.
	Datastore *DatastoreOpts
	// Networks dictates how this server will be attached to available networks.
	Networks []NetworkOpts
}

// ToInstanceCreateMap will render a JSON map.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]any, error) {
	if opts.Size > 300 || opts.Size < 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "instances.CreateOpts.Size"
		err.Value = opts.Size
		err.Info = "Size (GB) must be between 1-300"
		return nil, err
	}

	if opts.FlavorRef == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.CreateOpts.FlavorRef"}
	}

	instance := map[string]any{
		"flavorRef": opts.FlavorRef,
	}

	if opts.AvailabilityZone != "" {
		instance["availability_zone"] = opts.AvailabilityZone
	}

	if opts.Configuration != "" {
		instance["configuration"] = opts.Configuration
	}

	if opts.Name != "" {
		instance["name"] = opts.Name
	}
	if opts.Databases != nil {
		dbs, err := opts.Databases.ToDBCreateMap()
		if err != nil {
			return nil, err
		}
		instance["databases"] = dbs["databases"]
	}
	if opts.Users != nil {
		users, err := opts.Users.ToUserCreateMap()
		if err != nil {
			return nil, err
		}
		instance["users"] = users["users"]
	}
	if opts.Datastore != nil {
		datastore, err := opts.Datastore.ToMap()
		if err != nil {
			return nil, err
		}
		instance["datastore"] = datastore
	}

	if len(opts.Networks) > 0 {
		networks := make([]map[string]any, len(opts.Networks))
		for i, net := range opts.Networks {
			var err error
			networks[i], err = net.ToMap()
			if err != nil {
				return nil, err
			}
		}
		instance["nics"] = networks
	}

	volume := map[string]any{
		"size": opts.Size,
	}

	if opts.VolumeType != "" {
		volume["type"] = opts.VolumeType
	}

	instance["volume"] = volume

	return map[string]any{"instance": instance}, nil
}

// Create asynchronously provisions a new database instance. It requires the
// user to specify a flavor and a volume size. The API service then provisions
// the instance with the requested flavor and sets up a volume of the specified
// size, which is the storage for the database instance.
//
// Although this call only allows the creation of 1 instance per request, you
// can create an instance with multiple databases and users. The default
// binding for a MySQL instance is port 3306.
func Create(ctx context.Context, client gophercloud.Client, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List retrieves the status and information for all database instances.
func List(client gophercloud.Client) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves the status and information for a specified database instance.
func Get(ctx context.Context, client gophercloud.Client, id string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete permanently destroys the database instance.
func Delete(ctx context.Context, client gophercloud.Client, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// EnableRootUser enables the login from any host for the root user and
// provides the user with a generated root password.
func EnableRootUser(ctx context.Context, client gophercloud.Client, id string) (r EnableRootUserResult) {
	resp, err := client.Post(ctx, userRootURL(client, id), nil, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// IsRootEnabled checks an instance to see if root access is enabled. It returns
// True if root user is enabled for the specified database instance or False
// otherwise.
func IsRootEnabled(ctx context.Context, client gophercloud.Client, id string) (r IsRootEnabledResult) {
	resp, err := client.Get(ctx, userRootURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Restart will restart only the MySQL Instance. Restarting MySQL will
// erase any dynamic configuration settings that you have made within MySQL.
// The MySQL service will be unavailable until the instance restarts.
func Restart(ctx context.Context, client gophercloud.Client, id string) (r ActionResult) {
	b := map[string]any{"restart": struct{}{}}
	resp, err := client.Post(ctx, actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resize changes the memory size of the instance, assuming a valid
// flavorRef is provided. It will also restart the MySQL service.
func Resize(ctx context.Context, client gophercloud.Client, id, flavorRef string) (r ActionResult) {
	b := map[string]any{"resize": map[string]string{"flavorRef": flavorRef}}
	resp, err := client.Post(ctx, actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResizeVolume will resize the attached volume for an instance. It supports
// only increasing the volume size and does not support decreasing the size.
// The volume size is in gigabytes (GB) and must be an integer.
func ResizeVolume(ctx context.Context, client gophercloud.Client, id string, size int) (r ActionResult) {
	b := map[string]any{"resize": map[string]any{"volume": map[string]int{"size": size}}}
	resp, err := client.Post(ctx, actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AttachConfigurationGroup will attach configuration group to the instance
func AttachConfigurationGroup(ctx context.Context, client gophercloud.Client, instanceID string, configID string) (r ConfigurationResult) {
	b := map[string]any{"instance": map[string]any{"configuration": configID}}
	resp, err := client.Put(ctx, resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DetachConfigurationGroup will dettach configuration group from the instance
func DetachConfigurationGroup(ctx context.Context, client gophercloud.Client, instanceID string) (r ConfigurationResult) {
	b := map[string]any{"instance": map[string]any{}}
	resp, err := client.Put(ctx, resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
