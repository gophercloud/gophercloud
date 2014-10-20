package volumes

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

// Volume contains all the information associated with an OpenStack Volume.
type Volume struct {
	Status           string            `mapstructure:"status"`              // current status of the Volume
	Name             string            `mapstructure:"display_name"`        // display name
	Attachments      []string          `mapstructure:"attachments"`         // instances onto which the Volume is attached
	AvailabilityZone string            `mapstructure:"availability_zone"`   // logical group
	Bootable         string            `mapstructure:"bootable"`            // is the volume bootable
	CreatedAt        string            `mapstructure:"created_at"`          // date created
	Description      string            `mapstructure:"display_discription"` // display description
	VolumeType       string            `mapstructure:"volume_type"`         // see VolumeType object for more information
	SnapshotID       string            `mapstructure:"snapshot_id"`         // ID of the Snapshot from which the Volume was created
	SourceVolID      string            `mapstructure:"source_volid"`        // ID of the Volume from which the Volume was created
	Metadata         map[string]string `mapstructure:"metadata"`            // user-defined key-value pairs
	ID               string            `mapstructure:"id"`                  // unique identifier
	Size             int               `mapstructure:"size"`                // size of the Volume, in GB
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ListResult is a pagination.pager that is returned from a call to the List function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r ListResult) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractVolumes extracts and returns Volumes. It is used while iterating over a volumes.List call.
func ExtractVolumes(page pagination.Page) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Volumes, err
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Volume, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Volume *Volume `json:"volume"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Volume, err
}
