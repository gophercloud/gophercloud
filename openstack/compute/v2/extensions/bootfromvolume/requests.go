package bootfromvolume

import (
	"errors"
	"strconv"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"

	"github.com/racker/perigee"
)

// BlockDevice is a structure with options for booting a server instance
// from a volume. The volume may be created from an image, snapshot, or another
// volume.
type BlockDevice struct {
	// BootIndex [optional] is the boot index. It defaults to 0.
	BootIndex int `json:"boot_index"`

	// DeleteOnTermination [optional] specifies whether or not to delete the attached volume
	// when the server is deleted. Defaults to `false`.
	DeleteOnTermination bool `json:"delete_on_termination"`

	// DestinationType [optional] is the type that gets created. Possible values are "volume"
	// and "local".
	DestinationType string `json:"destination_type"`

	// SourceType [optional] must be one of: "volume", "snapshot", "image".
	SourceType string `json:"source_type"`

	// UUID [optional] is the unique identifier for the volume, snapshot, or image (see above)
	UUID string `json:"uuid"`

	// VolumeSize [optional] is the size of the volume to create (in gigabytes).
	VolumeSize int `json:"volume_size"`
}

// CreateOptsExt is a structure that extends the server `CreateOpts` structure
// by allowing for a block device mapping.
type CreateOptsExt struct {
	servers.CreateOptsBuilder
	BlockDevice BlockDevice `json:"block_device_mapping_v2,omitempty"`
}

// ToServerCreateMap adds the block device mapping option to the base server
// creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	var blockDevice BlockDevice
	if opts.BlockDevice == blockDevice {
		return base, nil
	}

	if opts.BlockDevice.SourceType != "volume" &&
		opts.BlockDevice.SourceType != "image" &&
		opts.BlockDevice.SourceType != "snapshot" &&
		opts.BlockDevice.SourceType != "" {
		return nil, errors.New("SourceType must be one of: volume, image, snapshot, [blank].")
	}

	serverMap := base["server"].(map[string]interface{})

	bd := make(map[string]interface{})
	bd["source_type"] = opts.BlockDevice.SourceType
	bd["boot_index"] = strconv.Itoa(opts.BlockDevice.BootIndex)
	bd["delete_on_termination"] = strconv.FormatBool(opts.BlockDevice.DeleteOnTermination)
	bd["volume_size"] = strconv.Itoa(opts.BlockDevice.VolumeSize)
	if opts.BlockDevice.UUID != "" {
		bd["uuid"] = opts.BlockDevice.UUID
	}
	if opts.BlockDevice.DestinationType != "" {
		bd["destination_type"] = opts.BlockDevice.DestinationType
	}

	serverMap["block_device_mapping_v2"] = []map[string]interface{}{bd}

	return base, nil
}

// Create requests the creation of a server from the given block device mapping.
func Create(client *gophercloud.ServiceClient, opts servers.CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200, 202},
		DumpReqJson: true,
	})
	return res
}
