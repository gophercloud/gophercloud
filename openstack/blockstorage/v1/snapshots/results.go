package snapshots

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

// Snapshot contains all the information associated with an OpenStack Snapshot.
type Snapshot struct {
	Status           string            `mapstructure:"status"`              // currect status of the Snapshot
	Name             string            `mapstructure:"display_name"`        // display name
	Attachments      []string          `mapstructure:"attachments"`         // instances onto which the Snapshot is attached
	AvailabilityZone string            `mapstructure:"availability_zone"`   // logical group
	Bootable         string            `mapstructure:"bootable"`            // is the Snapshot bootable
	CreatedAt        string            `mapstructure:"created_at"`          // date created
	Description      string            `mapstructure:"display_discription"` // display description
	VolumeType       string            `mapstructure:"volume_type"`         // see VolumeType object for more information
	SnapshotID       string            `mapstructure:"snapshot_id"`         // ID of the Snapshot from which this Snapshot was created
	SourceVolID      string            `mapstructure:"source_volid"`        // ID of the Volume from which this Snapshot was created
	Metadata         map[string]string `mapstructure:"metadata"`            // user-defined key-value pairs
	ID               string            `mapstructure:"id"`                  // unique identifier
	Size             int               `mapstructure:"size"`                // size of the Snapshot, in GB
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ListResult is a pagination.Pager that is returned from a call to the List function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Snapshots.
func (r ListResult) IsEmpty() (bool, error) {
	volumes, err := ExtractSnapshots(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractSnapshots extracts and returns Snapshots. It is used while iterating over a snapshots.List call.
func ExtractSnapshots(page pagination.Page) ([]Snapshot, error) {
	var response struct {
		Snapshots []Snapshot `json:"snapshots"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Snapshots, err
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

type commonResult struct {
	gophercloud.CommonResult
}

// Extract will get the Snapshot object out of the commonResult object.
func (r commonResult) Extract() (*Snapshot, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Snapshot *Snapshot `json:"snapshot"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("snapshots: Error decoding snapshots.commonResult: %v", err)
	}
	return res.Snapshot, nil
}
