package manageablevolumes

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

type ManageExistingResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the ManageExistingResult object.
func (r ManageExistingResult) Extract() (*volumes.Volume, error) {
	var s volumes.Volume
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a volume struct
func (r ManageExistingResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "volume")
}
