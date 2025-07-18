package servers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"regexp"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToServerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// ChangesSince is a time/date stamp for when the server last changed status.
	ChangesSince string `q:"changes-since"`

	// Image is the name of the image in URL format.
	Image string `q:"image"`

	// Flavor is the name of the flavor in URL format.
	Flavor string `q:"flavor"`

	// IP is a regular expression to match the IPv4 address of the server.
	IP string `q:"ip"`

	// This requires the client to be set to microversion 2.5 or later, unless
	// the user is an admin.
	// IP is a regular expression to match the IPv6 address of the server.
	IP6 string `q:"ip6"`

	// Name of the server as a string; can be queried with regular expressions.
	// Realize that ?name=bob returns both bob and bobb. If you need to match bob
	// only, you can use a regular expression matching the syntax of the
	// underlying database server implemented for Compute.
	Name string `q:"name"`

	// Status is the value of the status of the server so that you can filter on
	// "ACTIVE" for example.
	Status string `q:"status"`

	// Host is the name of the host as a string.
	Host string `q:"host"`

	// Marker is a UUID of the server at which you want to set a marker.
	Marker string `q:"marker"`

	// Limit is an integer value for the limit of values to return.
	Limit int `q:"limit"`

	// AllTenants is a bool to show all tenants.
	AllTenants bool `q:"all_tenants"`

	// TenantID lists servers for a particular tenant.
	// Setting "AllTenants = true" is required.
	TenantID string `q:"tenant_id"`

	// This requires the client to be set to microversion 2.83 or later, unless
	// the user is an admin.
	// UserID lists servers for a particular user.
	UserID string `q:"user_id"`

	// This requires the client to be set to microversion 2.26 or later.
	// Tags filters on specific server tags. All tags must be present for the server.
	Tags string `q:"tags"`

	// This requires the client to be set to microversion 2.26 or later.
	// TagsAny filters on specific server tags. At least one of the tags must be present for the server.
	TagsAny string `q:"tags-any"`

	// This requires the client to be set to microversion 2.26 or later.
	// NotTags filters on specific server tags. All tags must be absent for the server.
	NotTags string `q:"not-tags"`

	// This requires the client to be set to microversion 2.26 or later.
	// NotTagsAny filters on specific server tags. At least one of the tags must be absent for the server.
	NotTagsAny string `q:"not-tags-any"`

	// Display servers based on their availability zone (Admin only until microversion 2.82).
	AvailabilityZone string `q:"availability_zone"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListSimple makes a request against the API to list servers accessible to you.
func ListSimple(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToServerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// List makes a request against the API to list servers details accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToServerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// SchedulerHintOptsBuilder builds the scheduler hints into a serializable format.
type SchedulerHintOptsBuilder interface {
	ToSchedulerHintsMap() (map[string]any, error)
}

// SchedulerHintOpts represents a set of scheduling hints that are passed to the
// OpenStack scheduler.
type SchedulerHintOpts struct {
	// Group specifies a Server Group to place the instance in.
	Group string

	// DifferentHost will place the instance on a compute node that does not
	// host the given instances.
	DifferentHost []string

	// SameHost will place the instance on a compute node that hosts the given
	// instances.
	SameHost []string

	// Query is a conditional statement that results in compute nodes able to
	// host the instance.
	Query []any

	// TargetCell specifies a cell name where the instance will be placed.
	TargetCell string `json:"target_cell,omitempty"`

	// DifferentCell specifies cells names where an instance should not be placed.
	DifferentCell []string `json:"different_cell,omitempty"`

	// BuildNearHostIP specifies a subnet of compute nodes to host the instance.
	BuildNearHostIP string

	// AdditionalProperies are arbitrary key/values that are not validated by nova.
	AdditionalProperties map[string]any
}

// ToSchedulerHintsMap assembles a request body for scheduler hints.
func (opts SchedulerHintOpts) ToSchedulerHintsMap() (map[string]any, error) {
	sh := make(map[string]any)

	uuidRegex, _ := regexp.Compile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")

	if opts.Group != "" {
		if !uuidRegex.MatchString(opts.Group) {
			err := gophercloud.ErrInvalidInput{}
			err.Argument = "servers.schedulerhints.SchedulerHintOpts.Group"
			err.Value = opts.Group
			err.Info = "Group must be a UUID"
			return nil, err
		}
		sh["group"] = opts.Group
	}

	if len(opts.DifferentHost) > 0 {
		for _, diffHost := range opts.DifferentHost {
			if !uuidRegex.MatchString(diffHost) {
				err := gophercloud.ErrInvalidInput{}
				err.Argument = "servers.schedulerhints.SchedulerHintOpts.DifferentHost"
				err.Value = opts.DifferentHost
				err.Info = "The hosts must be in UUID format."
				return nil, err
			}
		}
		sh["different_host"] = opts.DifferentHost
	}

	if len(opts.SameHost) > 0 {
		for _, sameHost := range opts.SameHost {
			if !uuidRegex.MatchString(sameHost) {
				err := gophercloud.ErrInvalidInput{}
				err.Argument = "servers.schedulerhints.SchedulerHintOpts.SameHost"
				err.Value = opts.SameHost
				err.Info = "The hosts must be in UUID format."
				return nil, err
			}
		}
		sh["same_host"] = opts.SameHost
	}

	/*
		Query can be something simple like:
			 [">=", "$free_ram_mb", 1024]

			Or more complex like:
				['and',
					['>=', '$free_ram_mb', 1024],
					['>=', '$free_disk_mb', 200 * 1024]
				]

		Because of the possible complexity, just make sure the length is a minimum of 3.
	*/
	if len(opts.Query) > 0 {
		if len(opts.Query) < 3 {
			err := gophercloud.ErrInvalidInput{}
			err.Argument = "servers.schedulerhints.SchedulerHintOpts.Query"
			err.Value = opts.Query
			err.Info = "Must be a conditional statement in the format of [op,variable,value]"
			return nil, err
		}

		// The query needs to be sent as a marshalled string.
		b, err := json.Marshal(opts.Query)
		if err != nil {
			err := gophercloud.ErrInvalidInput{}
			err.Argument = "servers.schedulerhints.SchedulerHintOpts.Query"
			err.Value = opts.Query
			err.Info = "Must be a conditional statement in the format of [op,variable,value]"
			return nil, err
		}

		sh["query"] = string(b)
	}

	if opts.TargetCell != "" {
		sh["target_cell"] = opts.TargetCell
	}

	if len(opts.DifferentCell) > 0 {
		sh["different_cell"] = opts.DifferentCell
	}

	if opts.BuildNearHostIP != "" {
		if _, _, err := net.ParseCIDR(opts.BuildNearHostIP); err != nil {
			err := gophercloud.ErrInvalidInput{}
			err.Argument = "servers.schedulerhints.SchedulerHintOpts.BuildNearHostIP"
			err.Value = opts.BuildNearHostIP
			err.Info = "Must be a valid subnet in the form 192.168.1.1/24"
			return nil, err
		}
		ipParts := strings.Split(opts.BuildNearHostIP, "/")
		sh["build_near_host_ip"] = ipParts[0]
		sh["cidr"] = "/" + ipParts[1]
	}

	if opts.AdditionalProperties != nil {
		for k, v := range opts.AdditionalProperties {
			sh[k] = v
		}
	}

	if len(sh) == 0 {
		return sh, nil
	}

	return map[string]any{"os:scheduler_hints": sh}, nil
}

// Network is used within CreateOpts to control a new server's network
// attachments.
type Network struct {
	// UUID of a network to attach to the newly provisioned server.
	// Required unless Port is provided.
	UUID string

	// Port of a neutron network to attach to the newly provisioned server.
	// Required unless UUID is provided.
	Port string

	// FixedIP specifies a fixed IPv4 address to be used on this network.
	FixedIP string

	// Tag may contain an optional device role tag for the server's virtual
	// network interface. This can be used to identify network interfaces when
	// multiple networks are connected to one server.
	//
	// Requires microversion 2.32 through 2.36 or 2.42 or later.
	Tag string
}

type (
	// DestinationType represents the type of medium being used as the
	// destination of the bootable device.
	DestinationType string

	// SourceType represents the type of medium being used as the source of the
	// bootable device.
	SourceType string
)

const (
	// DestinationLocal DestinationType is for using an ephemeral disk as the
	// destination.
	DestinationLocal DestinationType = "local"

	// DestinationVolume DestinationType is for using a volume as the destination.
	DestinationVolume DestinationType = "volume"

	// SourceBlank SourceType is for a "blank" or empty source.
	SourceBlank SourceType = "blank"

	// SourceImage SourceType is for using images as the source of a block device.
	SourceImage SourceType = "image"

	// SourceSnapshot SourceType is for using a volume snapshot as the source of
	// a block device.
	SourceSnapshot SourceType = "snapshot"

	// SourceVolume SourceType is for using a volume as the source of block
	// device.
	SourceVolume SourceType = "volume"
)

// BlockDevice is a structure with options for creating block devices in a
// server. The block device may be created from an image, snapshot, new volume,
// or existing volume. The destination may be a new volume, existing volume
// which will be attached to the instance, ephemeral disk, or boot device.
type BlockDevice struct {
	// SourceType must be one of: "volume", "snapshot", "image", or "blank".
	SourceType SourceType `json:"source_type" required:"true"`

	// UUID is the unique identifier for the existing volume, snapshot, or
	// image (see above).
	UUID string `json:"uuid,omitempty"`

	// BootIndex is the boot index. It defaults to 0.
	BootIndex int `json:"boot_index"`

	// DeleteOnTermination specifies whether or not to delete the attached volume
	// when the server is deleted. Defaults to `false`.
	DeleteOnTermination bool `json:"delete_on_termination"`

	// DestinationType is the type that gets created. Possible values are "volume"
	// and "local".
	DestinationType DestinationType `json:"destination_type,omitempty"`

	// GuestFormat specifies the format of the block device.
	// Not specifying this will cause the device to be formatted to the default in Nova
	// which is currently vfat.
	// https://opendev.org/openstack/nova/src/commit/d0b459423dd81644e8d9382b6c87fabaa4f03ad4/nova/privsep/fs.py#L257
	GuestFormat string `json:"guest_format,omitempty"`

	// VolumeSize is the size of the volume to create (in gigabytes). This can be
	// omitted for existing volumes.
	VolumeSize int `json:"volume_size,omitempty"`

	// DeviceType specifies the device type of the block devices.
	// Examples of this are disk, cdrom, floppy, lun, etc.
	DeviceType string `json:"device_type,omitempty"`

	// DiskBus is the bus type of the block devices.
	// Examples of this are ide, usb, virtio, scsi, etc.
	DiskBus string `json:"disk_bus,omitempty"`

	// VolumeType is the volume type of the block device.
	// This requires Compute API microversion 2.67 or later.
	VolumeType string `json:"volume_type,omitempty"`

	// Tag is an arbitrary string that can be applied to a block device.
	// Information about the device tags can be obtained from the metadata API
	// and the config drive, allowing devices to be easily identified.
	// This requires Compute API microversion 2.42 or later.
	Tag string `json:"tag,omitempty"`
}

// Personality is an array of files that are injected into the server at launch.
type Personality []*File

// File is used within CreateOpts and RebuildOpts to inject a file into the
// server at launch.
// File implements the json.Marshaler interface, so when a Create or Rebuild
// operation is requested, json.Marshal will call File's MarshalJSON method.
type File struct {
	// Path of the file.
	Path string

	// Contents of the file. Maximum content size is 255 bytes.
	Contents []byte
}

// MarshalJSON marshals the escaped file, base64 encoding the contents.
func (f *File) MarshalJSON() ([]byte, error) {
	file := struct {
		Path     string `json:"path"`
		Contents string `json:"contents"`
	}{
		Path:     f.Path,
		Contents: base64.StdEncoding.EncodeToString(f.Contents),
	}
	return json.Marshal(file)
}

// DiskConfig represents one of the two possible settings for the DiskConfig
// option when creating, rebuilding, or resizing servers: Auto or Manual.
type DiskConfig string

const (
	// Auto builds a server with a single partition the size of the target flavor
	// disk and automatically adjusts the filesystem to fit the entire partition.
	// Auto may only be used with images and servers that use a single EXT3
	// partition.
	Auto DiskConfig = "AUTO"

	// Manual builds a server using whatever partition scheme and filesystem are
	// present in the source image. If the target flavor disk is larger, the
	// remaining space is left unpartitioned. This enables images to have non-EXT3
	// filesystems, multiple partitions, and so on, and enables you to manage the
	// disk configuration. It also results in slightly shorter boot times.
	Manual DiskConfig = "MANUAL"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerCreateMap() (map[string]any, error)
}

// CreateOpts specifies server creation parameters.
type CreateOpts struct {
	// Name is the name to assign to the newly launched server.
	Name string `json:"name" required:"true"`

	// ImageRef is the ID or full URL to the image that contains the
	// server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageRef string `json:"imageRef"`

	// FlavorRef is the ID or full URL to the flavor that describes the server's specs.
	FlavorRef string `json:"flavorRef"`

	// SecurityGroups lists the names of the security groups to which this server
	// should belong.
	SecurityGroups []string `json:"-"`

	// UserData contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you, if it isn't already.
	UserData []byte `json:"-"`

	// AvailabilityZone in which to launch the server.
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// Networks dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the
	// tenant.
	// Starting with microversion 2.37 networks can also be an "auto" or "none"
	// string.
	Networks any `json:"-"`

	// Metadata contains key-value pairs (up to 255 bytes each) to attach to the
	// server.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Personality includes files to inject into the server at launch.
	// Create will base64-encode file contents for you.
	Personality Personality `json:"personality,omitempty"`

	// ConfigDrive enables metadata injection through a configuration drive.
	ConfigDrive *bool `json:"config_drive,omitempty"`

	// AdminPass sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string `json:"adminPass,omitempty"`

	// AccessIPv4 specifies an IPv4 address for the instance.
	AccessIPv4 string `json:"accessIPv4,omitempty"`

	// AccessIPv6 specifies an IPv6 address for the instance.
	AccessIPv6 string `json:"accessIPv6,omitempty"`

	// Min specifies Minimum number of servers to launch.
	Min int `json:"min_count,omitempty"`

	// Max specifies Maximum number of servers to launch.
	Max int `json:"max_count,omitempty"`

	// Tags allows a server to be tagged with single-word metadata.
	// Requires microversion 2.52 or later.
	Tags []string `json:"tags,omitempty"`

	// (Available from 2.90) Hostname specifies the hostname to configure for the
	// instance in the metadata service. Starting with microversion 2.94, this can
	// be a Fully Qualified Domain Name (FQDN) of up to 255 characters in length.
	// If not set, OpenStack will derive the server's hostname from the Name field.
	Hostname string `json:"hostname,omitempty"`

	// BlockDevice describes the mapping of various block devices.
	BlockDevice []BlockDevice `json:"block_device_mapping_v2,omitempty"`

	// DiskConfig [optional] controls how the created server's disk is partitioned.
	DiskConfig DiskConfig `json:"OS-DCF:diskConfig,omitempty"`

	// KeyName is the name of the key pair.
	KeyName string `json:"key_name,omitempty"`

	// HypervisorHostname is the name of the hypervisor to which the server is scheduled.
	HypervisorHostname string `json:"hypervisor_hostname,omitempty"`
}

// ToServerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]any, error) {
	// We intentionally don't envelope the body here since we want to strip
	// some fields out and modify others
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.UserData != nil {
		var userData string
		if _, err := base64.StdEncoding.DecodeString(string(opts.UserData)); err != nil {
			userData = base64.StdEncoding.EncodeToString(opts.UserData)
		} else {
			userData = string(opts.UserData)
		}
		b["user_data"] = &userData
	}

	if len(opts.SecurityGroups) > 0 {
		securityGroups := make([]map[string]any, len(opts.SecurityGroups))
		for i, groupName := range opts.SecurityGroups {
			securityGroups[i] = map[string]any{"name": groupName}
		}
		b["security_groups"] = securityGroups
	}

	switch v := opts.Networks.(type) {
	case []Network:
		if len(v) > 0 {
			networks := make([]map[string]any, len(v))
			for i, net := range v {
				networks[i] = make(map[string]any)
				if net.UUID != "" {
					networks[i]["uuid"] = net.UUID
				}
				if net.Port != "" {
					networks[i]["port"] = net.Port
				}
				if net.FixedIP != "" {
					networks[i]["fixed_ip"] = net.FixedIP
				}
				if net.Tag != "" {
					networks[i]["tag"] = net.Tag
				}
			}
			b["networks"] = networks
		}
	case string:
		if v == "auto" || v == "none" {
			b["networks"] = v
		} else {
			return nil, fmt.Errorf(`networks must be a slice of Network struct or a string with "auto" or "none" values, current value is %q`, v)
		}
	}

	if opts.Min != 0 {
		b["min_count"] = opts.Min
	}

	if opts.Max != 0 {
		b["max_count"] = opts.Max
	}

	// Now we do our enveloping
	b = map[string]any{"server": b}

	return b, nil
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder, hintOpts SchedulerHintOptsBuilder) (r CreateResult) {
	b, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	if hintOpts != nil {
		sh, err := hintOpts.ToSchedulerHintsMap()
		if err != nil {
			r.Err = err
			return
		}
		maps.Copy(b, sh)
	}

	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests that a server previously provisioned be removed from your
// account.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ForceDelete forces the deletion of a server.
func ForceDelete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ActionResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"forceDelete": ""}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get requests details on a single server, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToServerUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing
// server.
type UpdateOpts struct {
	// Name changes the displayed name of the server.
	// The server host name will *not* change.
	// Server names are not constrained to be unique, even within the same tenant.
	Name *string `json:"name,omitempty"`

	// AccessIPv4 provides a new IPv4 address for the instance.
	AccessIPv4 *string `json:"accessIPv4,omitempty"`

	// AccessIPv6 provides a new IPv6 address for the instance.
	AccessIPv6 *string `json:"accessIPv6,omitempty"`

	// Hostname changes the hostname of the server.
	// Requires microversion 2.90 or later.
	// Note: This information is published via the metadata service and requires
	// application such as cloud-init to propagate it through to the instance.
	Hostname *string `json:"hostname,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToServerUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "server")
}

// Update requests that various attributes of the indicated server be changed.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToServerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ChangeAdminPassword alters the administrator or root password for a specified
// server.
func ChangeAdminPassword(ctx context.Context, client *gophercloud.ServiceClient, id, newPassword string) (r ActionResult) {
	b := map[string]any{
		"changePassword": map[string]string{
			"adminPass": newPassword,
		},
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RebootMethod describes the mechanisms by which a server reboot can be requested.
type RebootMethod string

// These constants determine how a server should be rebooted.
// See the Reboot() function for further details.
const (
	SoftReboot RebootMethod = "SOFT"
	HardReboot RebootMethod = "HARD"
	OSReboot                = SoftReboot
	PowerCycle              = HardReboot
)

// RebootOptsBuilder allows extensions to add additional parameters to the
// reboot request.
type RebootOptsBuilder interface {
	ToServerRebootMap() (map[string]any, error)
}

// RebootOpts provides options to the reboot request.
type RebootOpts struct {
	// Type is the type of reboot to perform on the server.
	Type RebootMethod `json:"type" required:"true"`
}

// ToServerRebootMap builds a body for the reboot request.
func (opts RebootOpts) ToServerRebootMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "reboot")
}

/*
Reboot requests that a given server reboot.

Two methods exist for rebooting a server:

HardReboot (aka PowerCycle) starts the server instance by physically cutting
power to the machine, or if a VM, terminating it at the hypervisor level.
It's done. Caput. Full stop.
Then, after a brief while, power is restored or the VM instance restarted.

SoftReboot (aka OSReboot) simply tells the OS to restart under its own
procedure.
E.g., in Linux, asking it to enter runlevel 6, or executing
"sudo shutdown -r now", or by asking Windows to rtart the machine.
*/
func Reboot(ctx context.Context, client *gophercloud.ServiceClient, id string, opts RebootOptsBuilder) (r ActionResult) {
	b, err := opts.ToServerRebootMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RebuildOptsBuilder allows extensions to provide additional parameters to the
// rebuild request.
type RebuildOptsBuilder interface {
	ToServerRebuildMap() (map[string]any, error)
}

// RebuildOpts represents the configuration options used in a server rebuild
// operation.
type RebuildOpts struct {
	// AdminPass is the server's admin password
	AdminPass string `json:"adminPass,omitempty"`

	// ImageRef is the ID of the image you want your server to be provisioned on.
	ImageRef string `json:"imageRef"`

	// Name to set the server to
	Name string `json:"name,omitempty"`

	// AccessIPv4 [optional] provides a new IPv4 address for the instance.
	AccessIPv4 string `json:"accessIPv4,omitempty"`

	// AccessIPv6 [optional] provides a new IPv6 address for the instance.
	AccessIPv6 string `json:"accessIPv6,omitempty"`

	// Metadata [optional] contains key-value pairs (up to 255 bytes each)
	// to attach to the server.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Personality [optional] includes files to inject into the server at launch.
	// Rebuild will base64-encode file contents for you.
	Personality Personality `json:"personality,omitempty"`

	// DiskConfig controls how the rebuilt server's disk is partitioned.
	DiskConfig DiskConfig `json:"OS-DCF:diskConfig,omitempty"`
}

// ToServerRebuildMap formats a RebuildOpts struct into a map for use in JSON
func (opts RebuildOpts) ToServerRebuildMap() (map[string]any, error) {
	if opts.DiskConfig != "" && opts.DiskConfig != Auto && opts.DiskConfig != Manual {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "servers.RebuildOpts.DiskConfig"
		err.Info = "Must be either diskconfig.Auto or diskconfig.Manual"
		return nil, err
	}

	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return map[string]any{"rebuild": b}, nil
}

// Rebuild will reprovision the server according to the configuration options
// provided in the RebuildOpts struct.
func Rebuild(ctx context.Context, client *gophercloud.ServiceClient, id string, opts RebuildOptsBuilder) (r RebuildResult) {
	b, err := opts.ToServerRebuildMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResizeOptsBuilder allows extensions to add additional parameters to the
// resize request.
type ResizeOptsBuilder interface {
	ToServerResizeMap() (map[string]any, error)
}

// ResizeOpts represents the configuration options used to control a Resize
// operation.
type ResizeOpts struct {
	// FlavorRef is the ID of the flavor you wish your server to become.
	FlavorRef string `json:"flavorRef" required:"true"`

	// DiskConfig [optional] controls how the resized server's disk is partitioned.
	DiskConfig DiskConfig `json:"OS-DCF:diskConfig,omitempty"`
}

// ToServerResizeMap formats a ResizeOpts as a map that can be used as a JSON
// request body for the Resize request.
func (opts ResizeOpts) ToServerResizeMap() (map[string]any, error) {
	if opts.DiskConfig != "" && opts.DiskConfig != Auto && opts.DiskConfig != Manual {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "servers.ResizeOpts.DiskConfig"
		err.Info = "Must be either diskconfig.Auto or diskconfig.Manual"
		return nil, err
	}

	return gophercloud.BuildRequestBody(opts, "resize")
}

// Resize instructs the provider to change the flavor of the server.
//
// Note that this implies rebuilding it.
//
// Unfortunately, one cannot pass rebuild parameters to the resize function.
// When the resize completes, the server will be in VERIFY_RESIZE state.
// While in this state, you can explore the use of the new server's
// configuration. If you like it, call ConfirmResize() to commit the resize
// permanently. Otherwise, call RevertResize() to restore the old configuration.
func Resize(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResizeOptsBuilder) (r ActionResult) {
	b, err := opts.ToServerResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ConfirmResize confirms a previous resize operation on a server.
// See Resize() for more details.
func ConfirmResize(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ActionResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"confirmResize": nil}, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201, 202, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RevertResize cancels a previous resize operation on a server.
// See Resize() for more details.
func RevertResize(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ActionResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"revertResize": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetMetadataOptsBuilder allows extensions to add additional parameters to
// the Reset request.
type ResetMetadataOptsBuilder interface {
	ToMetadataResetMap() (map[string]any, error)
}

// MetadataOpts is a map that contains key-value pairs.
type MetadataOpts map[string]string

// ToMetadataResetMap assembles a body for a Reset request based on the contents
// of a MetadataOpts.
func (opts MetadataOpts) ToMetadataResetMap() (map[string]any, error) {
	return map[string]any{"metadata": opts}, nil
}

// ToMetadataUpdateMap assembles a body for an Update request based on the
// contents of a MetadataOpts.
func (opts MetadataOpts) ToMetadataUpdateMap() (map[string]any, error) {
	return map[string]any{"metadata": opts}, nil
}

// ResetMetadata will create multiple new key-value pairs for the given server
// ID.
// Note: Using this operation will erase any already-existing metadata and
// create the new metadata provided. To keep any already-existing metadata,
// use the UpdateMetadatas or UpdateMetadata function.
func ResetMetadata(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResetMetadataOptsBuilder) (r ResetMetadataResult) {
	b, err := opts.ToMetadataResetMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, metadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Metadata requests all the metadata for the given server ID.
func Metadata(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetMetadataResult) {
	resp, err := client.Get(ctx, metadataURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMetadataOptsBuilder allows extensions to add additional parameters to
// the Create request.
type UpdateMetadataOptsBuilder interface {
	ToMetadataUpdateMap() (map[string]any, error)
}

// UpdateMetadata updates (or creates) all the metadata specified by opts for
// the given server ID. This operation does not affect already-existing metadata
// that is not specified by opts.
func UpdateMetadata(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateMetadataOptsBuilder) (r UpdateMetadataResult) {
	b, err := opts.ToMetadataUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, metadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MetadatumOptsBuilder allows extensions to add additional parameters to the
// Create request.
type MetadatumOptsBuilder interface {
	ToMetadatumCreateMap() (map[string]any, string, error)
}

// MetadatumOpts is a map of length one that contains a key-value pair.
type MetadatumOpts map[string]string

// ToMetadatumCreateMap assembles a body for a Create request based on the
// contents of a MetadataumOpts.
func (opts MetadatumOpts) ToMetadatumCreateMap() (map[string]any, string, error) {
	if len(opts) != 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "servers.MetadatumOpts"
		err.Info = "Must have 1 and only 1 key-value pair"
		return nil, "", err
	}
	metadatum := map[string]any{"meta": opts}
	var key string
	for k := range metadatum["meta"].(MetadatumOpts) {
		key = k
	}
	return metadatum, key, nil
}

// CreateMetadatum will create or update the key-value pair with the given key
// for the given server ID.
func CreateMetadatum(ctx context.Context, client *gophercloud.ServiceClient, id string, opts MetadatumOptsBuilder) (r CreateMetadatumResult) {
	b, key, err := opts.ToMetadatumCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, metadatumURL(client, id, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Metadatum requests the key-value pair with the given key for the given
// server ID.
func Metadatum(ctx context.Context, client *gophercloud.ServiceClient, id, key string) (r GetMetadatumResult) {
	resp, err := client.Get(ctx, metadatumURL(client, id, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMetadatum will delete the key-value pair with the given key for the
// given server ID.
func DeleteMetadatum(ctx context.Context, client *gophercloud.ServiceClient, id, key string) (r DeleteMetadatumResult) {
	resp, err := client.Delete(ctx, metadatumURL(client, id, key), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAddresses makes a request against the API to list the servers IP
// addresses.
func ListAddresses(client *gophercloud.ServiceClient, id string) pagination.Pager {
	return pagination.NewPager(client, listAddressesURL(client, id), func(r pagination.PageResult) pagination.Page {
		return AddressPage{pagination.SinglePageBase(r)}
	})
}

// ListAddressesByNetwork makes a request against the API to list the servers IP
// addresses for the given network.
func ListAddressesByNetwork(client *gophercloud.ServiceClient, id, network string) pagination.Pager {
	return pagination.NewPager(client, listAddressesByNetworkURL(client, id, network), func(r pagination.PageResult) pagination.Page {
		return NetworkAddressPage{pagination.SinglePageBase(r)}
	})
}

// CreateImageOptsBuilder allows extensions to add additional parameters to the
// CreateImage request.
type CreateImageOptsBuilder interface {
	ToServerCreateImageMap() (map[string]any, error)
}

// CreateImageOpts provides options to pass to the CreateImage request.
type CreateImageOpts struct {
	// Name of the image/snapshot.
	Name string `json:"name" required:"true"`

	// Metadata contains key-value pairs (up to 255 bytes each) to attach to
	// the created image.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ToServerCreateImageMap formats a CreateImageOpts structure into a request
// body.
func (opts CreateImageOpts) ToServerCreateImageMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "createImage")
}

// CreateImage makes a request against the nova API to schedule an image to be
// created of the server
func CreateImage(ctx context.Context, client *gophercloud.ServiceClient, id string, opts CreateImageOptsBuilder) (r CreateImageResult) {
	b, err := opts.ToServerCreateImageMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetPassword makes a request against the nova API to get the encrypted
// administrative password.
func GetPassword(ctx context.Context, client *gophercloud.ServiceClient, serverId string) (r GetPasswordResult) {
	resp, err := client.Get(ctx, passwordURL(client, serverId), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ShowConsoleOutputOptsBuilder is the interface types must satisfy in order to be
// used as ShowConsoleOutput options
type ShowConsoleOutputOptsBuilder interface {
	ToServerShowConsoleOutputMap() (map[string]any, error)
}

// ShowConsoleOutputOpts satisfies the ShowConsoleOutputOptsBuilder
type ShowConsoleOutputOpts struct {
	// The number of lines to fetch from the end of console log.
	// All lines will be returned if this is not specified.
	Length int `json:"length,omitempty"`
}

// ToServerShowConsoleOutputMap formats a ShowConsoleOutputOpts structure into a request body.
func (opts ShowConsoleOutputOpts) ToServerShowConsoleOutputMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "os-getConsoleOutput")
}

// ShowConsoleOutput makes a request against the nova API to get console log from the server
func ShowConsoleOutput(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ShowConsoleOutputOptsBuilder) (r ShowConsoleOutputResult) {
	b, err := opts.ToServerShowConsoleOutputMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// EvacuateOptsBuilder allows extensions to add additional parameters to the
// the Evacuate request.
type EvacuateOptsBuilder interface {
	ToEvacuateMap() (map[string]any, error)
}

// EvacuateOpts specifies Evacuate action parameters.
type EvacuateOpts struct {
	// The name of the host to which the server is evacuated
	Host string `json:"host,omitempty"`

	// Indicates whether server is on shared storage
	OnSharedStorage bool `json:"onSharedStorage"`

	// An administrative password to access the evacuated server
	AdminPass string `json:"adminPass,omitempty"`
}

// ToServerGroupCreateMap constructs a request body from CreateOpts.
func (opts EvacuateOpts) ToEvacuateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "evacuate")
}

// Evacuate will Evacuate a failed instance to another host.
func Evacuate(ctx context.Context, client *gophercloud.ServiceClient, id string, opts EvacuateOptsBuilder) (r EvacuateResult) {
	b, err := opts.ToEvacuateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// InjectNetworkInfo will inject the network info into a server
func InjectNetworkInfo(ctx context.Context, client *gophercloud.ServiceClient, id string) (r InjectNetworkResult) {
	b := map[string]any{
		"injectNetworkInfo": nil,
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Lock is the operation responsible for locking a Compute server.
func Lock(ctx context.Context, client *gophercloud.ServiceClient, id string) (r LockResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"lock": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unlock is the operation responsible for unlocking a Compute server.
func Unlock(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnlockResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"unlock": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Migrate will initiate a migration of the instance to another host.
func Migrate(ctx context.Context, client *gophercloud.ServiceClient, id string) (r MigrateResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"migrate": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// LiveMigrateOptsBuilder allows extensions to add additional parameters to the
// LiveMigrate request.
type LiveMigrateOptsBuilder interface {
	ToLiveMigrateMap() (map[string]any, error)
}

// LiveMigrateOpts specifies parameters of live migrate action.
type LiveMigrateOpts struct {
	// The host to which to migrate the server.
	// If this parameter is None, the scheduler chooses a host.
	Host *string `json:"host"`

	// Set to True to migrate local disks by using block migration.
	// If the source or destination host uses shared storage and you set
	// this value to True, the live migration fails.
	BlockMigration *bool `json:"block_migration,omitempty"`

	// Set to True to enable over commit when the destination host is checked
	// for available disk space. Set to False to disable over commit. This setting
	// affects only the libvirt virt driver.
	DiskOverCommit *bool `json:"disk_over_commit,omitempty"`
}

// ToLiveMigrateMap constructs a request body from LiveMigrateOpts.
func (opts LiveMigrateOpts) ToLiveMigrateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "os-migrateLive")
}

// LiveMigrate will initiate a live-migration (without rebooting) of the instance to another host.
func LiveMigrate(ctx context.Context, client *gophercloud.ServiceClient, id string, opts LiveMigrateOptsBuilder) (r MigrateResult) {
	b, err := opts.ToLiveMigrateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Pause is the operation responsible for pausing a Compute server.
func Pause(ctx context.Context, client *gophercloud.ServiceClient, id string) (r PauseResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"pause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unpause is the operation responsible for unpausing a Compute server.
func Unpause(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnpauseResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"unpause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RescueOptsBuilder is an interface that allows extensions to override the
// default structure of a Rescue request.
type RescueOptsBuilder interface {
	ToServerRescueMap() (map[string]any, error)
}

// RescueOpts represents the configuration options used to control a Rescue
// option.
type RescueOpts struct {
	// AdminPass is the desired administrative password for the instance in
	// RESCUE mode.
	// If it's left blank, the server will generate a password.
	AdminPass string `json:"adminPass,omitempty"`

	// RescueImageRef contains reference on an image that needs to be used as
	// rescue image.
	// If it's left blank, the server will be rescued with the default image.
	RescueImageRef string `json:"rescue_image_ref,omitempty"`
}

// ToServerRescueMap formats a RescueOpts as a map that can be used as a JSON
// request body for the Rescue request.
func (opts RescueOpts) ToServerRescueMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "rescue")
}

// Rescue instructs the provider to place the server into RESCUE mode.
func Rescue(ctx context.Context, client *gophercloud.ServiceClient, id string, opts RescueOptsBuilder) (r RescueResult) {
	b, err := opts.ToServerRescueMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unrescue instructs the provider to return the server from RESCUE mode.
func Unrescue(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnrescueResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"unrescue": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetNetwork will reset the network of a server
func ResetNetwork(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ResetNetworkResult) {
	b := map[string]any{
		"resetNetwork": nil,
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ServerState refers to the states usable in ResetState Action
type ServerState string

const (
	// StateActive returns the state of the server as active
	StateActive ServerState = "active"

	// StateError returns the state of the server as error
	StateError ServerState = "error"
)

// ResetState will reset the state of a server
func ResetState(ctx context.Context, client *gophercloud.ServiceClient, id string, state ServerState) (r ResetStateResult) {
	stateMap := map[string]any{"state": state}
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"os-resetState": stateMap}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Shelve is the operation responsible for shelving a Compute server.
func Shelve(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ShelveResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"shelve": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ShelveOffload is the operation responsible for Shelve-Offload a Compute server.
func ShelveOffload(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ShelveOffloadResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"shelveOffload": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UnshelveOptsBuilder allows extensions to add additional parameters to the
// Unshelve request.
type UnshelveOptsBuilder interface {
	ToUnshelveMap() (map[string]any, error)
}

// UnshelveOpts specifies parameters of shelve-offload action.
type UnshelveOpts struct {
	// Sets the availability zone to unshelve a server
	// Available only after nova 2.77
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts UnshelveOpts) ToUnshelveMap() (map[string]any, error) {
	// Key 'availabilty_zone' is required if the unshelve action is an object
	// i.e {"unshelve": {}} will be rejected
	b, err := gophercloud.BuildRequestBody(opts, "unshelve")
	if err != nil {
		return nil, err
	}

	if _, ok := b["unshelve"].(map[string]any)["availability_zone"]; !ok {
		b["unshelve"] = nil
	}

	return b, err
}

// Unshelve is the operation responsible for unshelve a Compute server.
func Unshelve(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UnshelveOptsBuilder) (r UnshelveResult) {
	b, err := opts.ToUnshelveMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Start is the operation responsible for starting a Compute server.
func Start(ctx context.Context, client *gophercloud.ServiceClient, id string) (r StartResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"os-start": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Stop is the operation responsible for stopping a Compute server.
func Stop(ctx context.Context, client *gophercloud.ServiceClient, id string) (r StopResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"os-stop": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Suspend is the operation responsible for suspending a Compute server.
func Suspend(ctx context.Context, client *gophercloud.ServiceClient, id string) (r SuspendResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"suspend": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resume is the operation responsible for resuming a Compute server.
func Resume(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ResumeResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"resume": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
