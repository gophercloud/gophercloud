package volumeTypes

import (
//"fmt"

//"github.com/rackspace/gophercloud/pagination"

//"github.com/mitchellh/mapstructure"
)

type VolumeType struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs" mapstructure:"extra_specs"`
	ID         string                 `json:"id" mapstructure:"id"`
	Name       string                 `json:"name" mapstructure:"name"`
}
