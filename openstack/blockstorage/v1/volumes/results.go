package volumes

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

type Volume struct {
	Status           string            `mapstructure:"status"`
	Name             string            `mapstructure:"display_name"`
	Attachments      []string          `mapstructure:"attachments"`
	AvailabilityZone string            `mapstructure:"availability_zone"`
	Bootable         string            `mapstructure:"bootable"`
	CreatedAt        string            `mapstructure:"created_at"`
	Description      string            `mapstructure:"display_discription"`
	VolumeType       string            `mapstructure:"volume_type"`
	SnapshotID       string            `mapstructure:"snapshot_id"`
	SourceVolID      string            `mapstructure:"source_volid"`
	Metadata         map[string]string `mapstructure:"metadata"`
	ID               string            `mapstructure:"id"`
	Size             int               `mapstructure:"size"`
}

// ListResult is a *http.Response that is returned from a call to the List function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no container names.
func (r ListResult) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractVolumes extracts and returns the Volumes from a 'List' request.
func ExtractVolumes(page pagination.Page) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Volumes, err
}

type commonResult struct {
	gophercloud.CommonResult
}

// ExtractVolume extracts and returns the Volume from a 'Get' request.
func (r commonResult) Extract() (*Volume, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Volume *Volume `json:"volume"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("volumes: Error decoding volumes.commonResult: %v", err)
	}
	return res.Volume, nil
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}
type UpdateResult struct {
	commonResult
}

type DeleteResult commonResult
