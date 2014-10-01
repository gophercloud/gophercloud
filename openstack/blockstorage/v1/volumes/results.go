package volumes

import (
	"fmt"

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

// ListOpts holds options for listing volumes. It is passed to the volumes.List function.
type ListOpts struct {
	// AllTenants is an admin-only option. Set it to true to see a tenant volumes.
	AllTenants bool
	// List only volumes that contain Metadata.
	Metadata map[string]string
	// List only volumes that have Name as the display name.
	Name string
	// List only volumes that have a status of Status.
	Status string
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

type GetResult struct {
	err error
	r   map[string]interface{}
}

// ExtractVolume extracts and returns the Volume from a 'Get' request.
func (gr GetResult) ExtractVolume() (*Volume, error) {
	if gr.err != nil {
		return nil, gr.err
	}

	var response struct {
		Volume *Volume `json:"volume"`
	}

	err := mapstructure.Decode(gr.r, &response)
	if err != nil {
		return nil, fmt.Errorf("volumes: Error decoding volumes.GetResult: %v", err)
	}
	return response.Volume, nil
}
