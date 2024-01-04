package instances

import (
	"github.com/gophercloud/gophercloud"
	db "github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder is the top-level interface for create options.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// DatastoreOpts represents the configuration for how an instance stores data.
type DatastoreOpts struct {
	Version string `json:"version"`
	Type    string `json:"type"`
}

// ToMap converts a DatastoreOpts to a map[string]string (for a request body)
func (opts DatastoreOpts) ToMap() (map[string]interface{}, error) {
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
func (opts NetworkOpts) ToMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AccessOpts is used within CreateOpts to define how the database service is exposed.
type AccessOpts struct {
	// Specifies whether the database service is exposed to the public
	IsPublic bool `json:"is_public,omitempty"`
	// A list of IPv4, IPv6 or mix of both CIDRs that restrict access to the database service
	// 0.0.0.0/0 is used by default if this parameter is not provided.
	AllowedCidrs []string `json:"allowed_cidrs,omitempty"`
}

// ToMap converts a AccessOpts to a map[string]string (for a request body)
func (opts AccessOpts) ToMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// CreateOpts is the struct responsible for configuring a new database instance.
type CreateOpts struct {
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
	// Specifies the availability zone of the instance. Optional
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	// ID of the configuration group that you want to attach to the instance.
	// Optional.
	Configuration string `json:"configuration,omitempty"`
	// Specifies how the database service is exposed
	Access *AccessOpts `json:"access,omitempty"`
	// The scheduler hint when creating underlying Nova instances.
	// Valide values are: affinity, anti-affinity.
	Locality string `json:"locality,omitempty"`
}

// ToInstanceCreateMap will render a JSON map.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
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

	instance := map[string]interface{}{
		"flavorRef": opts.FlavorRef,
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
		networks := make([]map[string]interface{}, len(opts.Networks))
		for i, net := range opts.Networks {
			var err error
			networks[i], err = net.ToMap()
			if err != nil {
				return nil, err
			}
		}
		instance["nics"] = networks
	}

	volume := map[string]interface{}{
		"size": opts.Size,
	}

	if opts.VolumeType != "" {
		volume["type"] = opts.VolumeType
	}

	instance["volume"] = volume

	if opts.AvailabilityZone != "" {
		instance["availability_zone"] = opts.AvailabilityZone
	}

	if opts.Configuration != "" {
		instance["configuration"] = opts.Configuration
	}

	if opts.Access != nil {
		access, err := opts.Access.ToMap()
		if err != nil {
			return nil, err
		}
		instance["access"] = access
	}

	if opts.Locality != "" {
		instance["locality"] = opts.Locality
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
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List retrieves the status and information for all database instances.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves the status and information for a specified database instance.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(resourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete permanently destroys the database instance.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(resourceURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// EnableRootUser enables the login from any host for the root user and
// provides the user with a generated root password.
func EnableRootUser(client *gophercloud.ServiceClient, id string) (r EnableRootUserResult) {
	resp, err := client.Post(userRootURL(client, id), nil, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DisableRootUser disables the login for the root user
func DisableRootUser(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(userRootURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// IsRootEnabled checks an instance to see if root access is enabled. It returns
// True if root user is enabled for the specified database instance or False
// otherwise.
func IsRootEnabled(client *gophercloud.ServiceClient, id string) (r IsRootEnabledResult) {
	resp, err := client.Get(userRootURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Restart will restart only the MySQL Instance. Restarting MySQL will
// erase any dynamic configuration settings that you have made within MySQL.
// The MySQL service will be unavailable until the instance restarts.
func Restart(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"restart": struct{}{}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resize changes the memory size of the instance, assuming a valid
// flavorRef is provided. It will also restart the MySQL service.
func Resize(client *gophercloud.ServiceClient, id, flavorRef string) (r ActionResult) {
	b := map[string]interface{}{"resize": map[string]string{"flavorRef": flavorRef}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResizeVolume will resize the attached volume for an instance. It supports
// only increasing the volume size and does not support decreasing the size.
// The volume size is in gigabytes (GB) and must be an integer.
func ResizeVolume(client *gophercloud.ServiceClient, id string, size int) (r ActionResult) {
	b := map[string]interface{}{"resize": map[string]interface{}{"volume": map[string]int{"size": size}}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AttachConfigurationGroup will attach configuration group to the instance
func AttachConfigurationGroup(client *gophercloud.ServiceClient, instanceID string, configID string) (r ConfigurationResult) {
	b := map[string]interface{}{"instance": map[string]interface{}{"configuration": configID}}
	resp, err := client.Put(resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DetachConfigurationGroup will dettach configuration group from the instance
func DetachConfigurationGroup(client *gophercloud.ServiceClient, instanceID string) (r ConfigurationResult) {
	b := map[string]interface{}{"instance": map[string]interface{}{"configuration": nil}}
	resp, err := client.Put(resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DetachReplica will detach a replica from its replication source.
// Detaching replica from the source will make the replica a standalone instance
func DetachReplica(client *gophercloud.ServiceClient, instanceID, replicaOf string) (r ReplicaResult) {
	b := map[string]interface{}{"instance": map[string]interface{}{"replica_of": replicaOf}}
	resp, err := client.Put(resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateInstanceAccessbility will update the accessbility of the instance
//   - Specifies if the instance should be exposed to public or not. Default is private
//   - Specifies the list of CIDRs that are allowd to access the database service.
//     Not set means allowing everything
func UpdateInstanceAccessbility(client *gophercloud.ServiceClient, instanceID string, opts AccessOpts) (r ReplicaResult) {
	access, err := opts.ToMap()
	if err != nil {
		r.Err = err
		return
	}
	b := map[string]interface{}{"instance": map[string]interface{}{"access": access}}
	resp, err := client.Put(resourceURL(client, instanceID), &b, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Promote will promote a replica to master
// Once this command is executed, the status of all the instances will be set to
// PROMOTE and Trove will work its magic until all of them to come back to HEALTHY.
func Promote(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"promote_to_replica_source": struct{}{}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Eject will eject the current master and force a re-eclection for the new master.
// Once this command is executed, the status of all the instances will be set to
// EJECT and Trove will work its magic until all of them to come back to HEALTHY.
func Eject(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"eject_replica_source": struct{}{}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Reset will set instance service status to ERROR and clear the current task status.
// Mark any running backup operations as FAILED.
func Reset(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"reset_status": struct{}{}}
	resp, err := client.Post(actionURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Stop will stop database service inside an instance.
// Admin only API.
func Stop(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"stop": struct{}{}}
	resp, err := client.Post(actionadminURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Reboot will stop database service inside an instance.
// Admin only API.
func Reboot(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"reboot": struct{}{}}
	resp, err := client.Post(actionadminURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetTask will reset task status of an instance, mark any running backup operations as FAILED
// Admin only API.
func ResetTask(client *gophercloud.ServiceClient, id string) (r ActionResult) {
	b := map[string]interface{}{"reset-task-status": struct{}{}}
	resp, err := client.Post(actionadminURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Rebuild will rebuild the Nova server’s operating system for the database instance.
// Communication with the end users is needed as the database service goes offline during the process
// User’s data in the database is not affected.
// Admin only API.
func Rebuild(client *gophercloud.ServiceClient, id, imageId string) (r ActionResult) {
	b := map[string]interface{}{"rebuild": map[string]interface{}{"image_id": imageId}}
	resp, err := client.Post(actionadminURL(client, id), &b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ReplicateOpts is the struct responsible for configuring a new database instance.
type ReplicateOpts struct {
	// Name of the instance to create. The length of the name is limited to
	// 255 characters and any characters are permitted. Optional.
	Name string `json:"name"`
	// Networks dictates how this server will be attached to available networks.
	Networks []NetworkOpts `json:"nics,omitempty"`
	// Specifies the unique ID or name of an existing instance to replicate from.
	ReplicaOf string `json:"replica_of"`
	// Number of replicas to create (defaults to 1). Optional.
	ReplicaCount int `json:"replica_count,omitempty"`
	// Specifies how the database service is exposed
	Access *AccessOpts `json:"access,omitempty"`
}

// ToMap converts a ReplicateOpts to a map[string]string (for a request body)
func (opts ReplicateOpts) ToMap() (map[string]interface{}, error) {
	if opts.ReplicaOf == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.ReplicateOpts.ReplicaOf"}
	}
	if opts.ReplicaCount > 3 || opts.ReplicaCount < 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "instances.ReplicateOpts.ReplicaCount"
		err.Value = opts.ReplicaCount
		err.Info = "Replica Count must be between 1-3"
		return nil, err
	}

	if opts.ReplicaOf == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.ReplicateOpts.ReplicaOf"}
	}

	instance := map[string]interface{}{
		"replica_of": opts.ReplicaOf,
	}

	if opts.Name != "" {
		instance["name"] = opts.Name
	}

	if len(opts.Networks) > 0 {
		networks := make([]map[string]interface{}, len(opts.Networks))
		for i, net := range opts.Networks {
			var err error
			networks[i], err = net.ToMap()
			if err != nil {
				return nil, err
			}
		}
		instance["nics"] = networks
	}

	instance["replica_count"] = opts.ReplicaCount

	if opts.Access != nil {
		access, err := opts.Access.ToMap()
		if err != nil {
			return nil, err
		}
		instance["access"] = access
	}

	return map[string]interface{}{"instance": instance}, nil
}

// ReplicateOptsBuilder is a top-level interface which renders a JSON map.
type ReplicateOptsBuilder interface {
	ToMap() (map[string]interface{}, error)
}

// Replicate an instance from an existing database instance.
func Replicate(client *gophercloud.ServiceClient, opts ReplicateOptsBuilder) (r CreateResult) {
	b, err := opts.ToMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RestoreOpts is the struct responsible for configuring a new database instance.
type RestoreOpts struct {
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
	Name string `json:"name"`
	// Options to configure the type of datastore the instance will use. This is
	// optional, and if excluded will default to MySQL.
	Datastore *DatastoreOpts
	// Networks dictates how this server will be attached to available networks.
	Networks []NetworkOpts
	// The backup id of the backup used from which a new instance is created.
	BackupRef string `json:"backupRef"`
	// Number of replicas to create (defaults to 1). Optional.
	ReplicaCount int `json:"replica_count,omitempty"`
	// Specifies how the database service is exposed
	Access *AccessOpts `json:"access,omitempty"`
}

// ToMap converts a RestoreOpts to a map[string]string (for a request body)
func (opts RestoreOpts) ToMap() (map[string]interface{}, error) {
	if opts.Size > 300 || opts.Size < 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "instances.RestoreOpts.Size"
		err.Value = opts.Size
		err.Info = "Size (GB) must be between 1-300"
		return nil, err
	}

	if opts.ReplicaCount > 3 || opts.ReplicaCount < 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "instances.RestoreOpts.ReplicaCount"
		err.Value = opts.ReplicaCount
		err.Info = "Replica Count must be between 1-3"
		return nil, err
	}

	if opts.BackupRef == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.RestoreOpts.BackupRef"}
	}

	if opts.FlavorRef == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.RestoreOpts.FlavorRef"}
	}

	restorePoint := map[string]interface{}{
		"backupRef": opts.BackupRef,
	}

	instance := map[string]interface{}{
		"flavorRef":    opts.FlavorRef,
		"restorePoint": restorePoint,
	}

	if opts.Name != "" {
		instance["name"] = opts.Name
	}
	if opts.Datastore != nil {
		datastore, err := opts.Datastore.ToMap()
		if err != nil {
			return nil, err
		}
		instance["datastore"] = datastore
	}

	if len(opts.Networks) > 0 {
		networks := make([]map[string]interface{}, len(opts.Networks))
		for i, net := range opts.Networks {
			var err error
			networks[i], err = net.ToMap()
			if err != nil {
				return nil, err
			}
		}
		instance["nics"] = networks
	}

	volume := map[string]interface{}{
		"size": opts.Size,
	}

	if opts.VolumeType != "" {
		volume["type"] = opts.VolumeType
	}

	instance["volume"] = volume
	instance["replica_count"] = opts.ReplicaCount

	if opts.Access != nil {
		access, err := opts.Access.ToMap()
		if err != nil {
			return nil, err
		}
		instance["access"] = access
	}

	return map[string]interface{}{"instance": instance}, nil
}

// RestoreOptsBuilder is a top-level interface which renders a JSON map.
type RestoreOptsBuilder interface {
	ToMap() (map[string]interface{}, error)
}

// Restore an instance from a backup database.
func Restore(client *gophercloud.ServiceClient, opts RestoreOptsBuilder) (r CreateResult) {
	b, err := opts.ToMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
