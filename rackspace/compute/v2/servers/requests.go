package servers

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/diskconfig"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// CreateOpts specifies all of the options that Rackspace accepts in its Create request, including
// the union of all extensions that Rackspace supports.
type CreateOpts struct {
	// Name [required] is the name to assign to the newly launched server.
	Name string

	// ImageRef [required] is the ID or full URL to the image that contains the server's OS and initial state.
	// Optional if using the boot-from-volume extension.
	ImageRef string

	// FlavorRef [required] is the ID or full URL to the flavor that describes the server's specs.
	FlavorRef string

	// SecurityGroups [optional] lists the names of the security groups to which this server should belong.
	SecurityGroups []string

	// UserData [optional] contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you.
	UserData []byte

	// AvailabilityZone [optional] in which to launch the server.
	AvailabilityZone string

	// Networks [optional] dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the tenant.
	Networks []os.Network

	// Metadata [optional] contains key-value pairs (up to 255 bytes each) to attach to the server.
	Metadata map[string]string

	// Personality [optional] includes the path and contents of a file to inject into the server at launch.
	// The maximum size of the file is 255 bytes (decoded).
	Personality []byte

	// ConfigDrive [optional] enables metadata injection through a configuration drive.
	ConfigDrive bool

	// Rackspace-specific extensions begin here.

	// KeyPair [optional] specifies the name of the SSH KeyPair to be injected into the newly launched
	// server. See the "keypairs" extension in OpenStack compute v2.
	KeyPair string

	// DiskConfig [optional] controls how the created server's disk is partitioned. See the "diskconfig"
	// extension in OpenStack compute v2.
	DiskConfig diskconfig.DiskConfig
}

// ToServerCreateMap constructs a request body using all of the OpenStack extensions that are
// active on Rackspace.
func (opts CreateOpts) ToServerCreateMap() map[string]interface{} {
	base := os.CreateOpts{
		Name:             opts.Name,
		ImageRef:         opts.ImageRef,
		FlavorRef:        opts.FlavorRef,
		SecurityGroups:   opts.SecurityGroups,
		UserData:         opts.UserData,
		AvailabilityZone: opts.AvailabilityZone,
		Networks:         opts.Networks,
		Metadata:         opts.Metadata,
		Personality:      opts.Personality,
		ConfigDrive:      opts.ConfigDrive,
	}

	drive := diskconfig.CreateOptsExt{
		CreateOptsBuilder: base,
		DiskConfig:        opts.DiskConfig,
	}

	result := drive.ToServerCreateMap()

	// key_name doesn't actually come from the extension (or at least isn't documented there) so
	// we need to add it manually.
	serverMap := result["server"].(map[string]interface{})
	serverMap["key_name"] = opts.KeyPair

	return result
}
