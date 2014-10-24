package bootfromvolume

import (
	"errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type CreateOptsExt struct {
	servers.CreateOptsBuilder
	BDM BlockDeviceMapping
}

// ToServerCreateMap adds the block device mapping option to the base server
// creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	if opts.BDM.SourceType != "volume" && opts.BDM.SourceType != "image" && opts.BDM.SourceType != "snapshot" {
		return nil, errors.New("SourceType must be one of: volume, image, snapshot.")
	}

	if opts.BDM.UUID == "" {
		return nil, errors.New("Required field UUID not set.")
	}

	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	serverMap := base["server"].(map[string]interface{})
	serverMap["block_device_mapping_v2"] = opts.BDM

	return base, nil
}
