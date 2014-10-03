package volumetypes

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

type commonResult struct {
	gophercloud.CommonResult
}

func (r commonResult) Extract() (*VolumeType, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		VolumeType *VolumeType `json:"volume_type" mapstructure:"volume_type"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Volume Type: %v", err)
	}

	return res.VolumeType, nil
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}
