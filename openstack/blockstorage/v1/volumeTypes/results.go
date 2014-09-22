package volumeTypes

import (
	"fmt"

	//"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

type VolumeType struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs" mapstructure:"extra_specs"`
	ID         string                 `json:"id" mapstructure:"id"`
	Name       string                 `json:"name" mapstructure:"name"`
}

type GetResult map[string]interface{}

func ExtractVolumeType(gr GetResult) (*VolumeType, error) {
	var response struct {
		VolumeType *VolumeType `mapstructure:"volume_type"`
	}
	err := mapstructure.Decode(gr, &response)
	if err != nil {
		return nil, fmt.Errorf("volumeTypes: Error decoding volumeTypes.GetResult: %v", err)
	}
	return response.VolumeType, nil
}
