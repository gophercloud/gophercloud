package bootfromvolume

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// SourceType represents the type of medium being used to create the volume.
type SourceType string

const (
	// Volume SourceType
	Volume SourceType = "volume"
	// Snapshot SourceType
	Snapshot SourceType = "snapshot"
	// Image SourceType
	Image SourceType = "image"
	// Blank SourceType
	Blank SourceType = "blank"
)

// BlockDevice is a structure with options for booting a server instance
// from a volume. The volume may be created from an image, snapshot, or another
// volume.
type BlockDevice struct {
	// SourceType must be one of: "volume", "snapshot", "image".
	SourceType SourceType `json:"source_type" required:"true"`
	// UUID is the unique identifier for the volume, snapshot, or image (see above)
	UUID string `json:"uuid,omitempty"`
	// BootIndex is the boot index. It defaults to 0.
	BootIndex int `json:"boot_index"`
	// DeleteOnTermination specifies whether or not to delete the attached volume
	// when the server is deleted. Defaults to `false`.
	DeleteOnTermination bool `json:"delete_on_termination"`
	// DestinationType is the type that gets created. Possible values are "volume"
	// and "local".
	DestinationType string `json:"destination_type,omitempty"`
	// GuestFormat specifies the format of the block device.
	GuestFormat string `json:"guest_format,omitempty"`
	// VolumeSize is the size of the volume to create (in gigabytes).
	VolumeSize int `json:"volume_size"`
}

// CreateOptsExt is a structure that extends the server `CreateOpts` structure
// by allowing for a block device mapping.
type CreateOptsExt struct {
	servers.CreateOptsBuilder
	BlockDevice []BlockDevice `json:"block_device_mapping_v2,omitempty"`
}

// ToServerCreateMap adds the block device mapping option to the base server
// creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	if len(opts.BlockDevice) == 0 {
		err := gophercloud.ErrMissingInput{}
		err.Argument = "bootfromvolume.CreateOptsExt.BlockDevice"
		return nil, err
	}

	serverMap := base["server"].(map[string]interface{})

	blockDevice := make([]map[string]interface{}, len(opts.BlockDevice))

	for i, bd := range opts.BlockDevice {
		b, err := gophercloud.BuildRequestBody(bd, "")
		if err != nil {
			return nil, err
		}
		blockDevice[i] = b
	}
	serverMap["block_device_mapping_v2"] = blockDevice

	return base, nil
}

// Create requests the creation of a server from the given block device mapping.
func Create(client *gophercloud.ServiceClient, opts servers.CreateOptsBuilder) (r servers.CreateResult) {
	b, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}
