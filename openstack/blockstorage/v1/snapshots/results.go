package snapshots

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

type Snapshot struct {
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
	volumes, err := ExtractSnapshots(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractSnapshots extracts and returns the Volumes from a 'List' request.
func ExtractSnapshots(page pagination.Page) ([]Snapshot, error) {
	var response struct {
		Snapshots []Snapshot `json:"snapshots"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Snapshots, err
}

type commonResult struct {
	gophercloud.CommonResult
}

// Extract returns a pointer to the Snapshot from a commonResult.Resp.
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
