package bootfromvolume

import (
  "errors"

  "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type CreateOptsExt struct {
	servers.CreateOptsBuilder
	BlockDeviceMapping BlockDeviceMapping
}

// ToServerCreateMap adds the block device mapping option to the base server
// creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
  if opts.SourceType != "volume" && opts.SourceType != "image" && opts.SourceType != "snapshot" {
    return nil, errors.New("SourceType must be one of: volume, image, snapshot.")
  }

  if opts.UUID == "" {
    return nil, errors.New("Required field UUID not set.")
  }

  base := opts.CreateOptsBuilder.ToServerCreateMap()

  serverMap := base["server"].(map[string]interface{})
  serverMap["block_device_mapping_v2"] = opts.BlockDeviceMapping

  return base
}
