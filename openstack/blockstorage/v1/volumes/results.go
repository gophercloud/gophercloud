package volumes

import (
	"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

type Volume struct {
	Status           string
	Name             string
	Attachments      []string
	AvailabilityZone string
	Bootable         bool
	CreatedAt        string
	Description      string
	VolumeType       string
	SnapshotID       string
	SourceVolID      string
	Metadata         map[string]string
	Id               string
	Size             int
}

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

func ExtractVolumes(page pagination.Page) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Volumes, err
}
