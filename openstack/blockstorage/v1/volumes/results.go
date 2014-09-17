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

type VolumePage struct {
	pagination.LinkedPageBase
}

func (p VolumPage) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(p)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

func ExtractVolumes(page pagination.page) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	err := mapstructure.Decode(page.(VolumePage).Body, &response)
	return response.Volumes, err
}

type CreateOpts map[string]interface{}
type ListOpts map[string]bool
type GetOpts map[string]string
type DeleteOpts map[string]string
