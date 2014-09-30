package volumeTypes

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type VolumeType struct {
	ExtraSpecs map[string]interface{} `json:"extra_specs" mapstructure:"extra_specs"`
	ID         string                 `json:"id" mapstructure:"id"`
	Name       string                 `json:"name" mapstructure:"name"`
}

// ListOpts holds options for listing volumes. It is passed to the volumes.List function.
type ListOpts struct {
}

// ListResult is a *http.Response that is returned from a call to the List function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no container names.
func (r ListResult) IsEmpty() (bool, error) {
	volumeTypes, err := ExtractVolumeTypes(r)
	if err != nil {
		return true, err
	}
	return len(volumeTypes) == 0, nil
}

// ExtractVolumeTypes extracts and returns the Volumes from a 'List' request.
func ExtractVolumeTypes(page pagination.Page) ([]VolumeType, error) {
	var response struct {
		VolumeTypes []VolumeType `mapstructure:"volume_types"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.VolumeTypes, err
}

type GetResult struct {
	gophercloud.CommonResult
}

func (gr GetResult) ExtractVolumeType() (*VolumeType, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}
	var response struct {
		VolumeType *VolumeType `mapstructure:"volume_type"`
	}
	err := mapstructure.Decode(gr, &response)
	if err != nil {
		return nil, fmt.Errorf("volumeTypes: Error decoding volumeTypes.GetResult: %v", err)
	}
	return response.VolumeType, nil
}
