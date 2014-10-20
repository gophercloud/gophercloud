package snapshots

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

// Status is the type used to represent a snapshot's status
type Status string

// Constants to use for supported statuses
const (
	Creating    Status = "CREATING"
	Available   Status = "AVAILABLE"
	Deleting    Status = "DELETING"
	Error       Status = "ERROR"
	DeleteError Status = "ERROR_DELETING"
)

// Snapshot is the Rackspace representation of an external block storage device.
type Snapshot struct {
	// The timestamp when this snapshot was created.
	CreatedAt string `mapstructure:"created_at"`

	// The human-readable description for this snapshot.
	Description string `mapstructure:"display_description"`

	// The human-readable name for this snapshot.
	Name string `mapstructure:"display_name"`

	// The UUID for this snapshot.
	ID string `mapstructure:"id"`

	// The random metadata associated with this snapshot. Note: unlike standard
	// OpenStack snapshots, this cannot actually be set.
	Metadata map[string]string `mapstructure:"metadata"`

	// Indicates the current progress of the snapshot's backup procedure.
	Progress string `mapstructure:"os-extended-snapshot-attributes:progress"`

	// The project ID.
	ProjectID string `mapstructure:"os-extended-snapshot-attributes:project_id"`

	// The size of the volume which this snapshot backs up.
	Size int `mapstructure:"size"`

	// The status of the snapshot.
	Status Status `mapstructure:"status"`

	// The ID of the volume which this snapshot seeks to back up.
	VolumeID string `mapstructure:"volume_id"`
}

type commonResult struct {
	gophercloud.CommonResult
}

// CreateResult represents the result of a create operation
type CreateResult struct {
	Common os.CreateResult
	commonResult
}

// GetResult represents the result of a get operation
type GetResult struct {
	Common os.GetResult
	commonResult
}
