package snapshots

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	os "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSnapshotCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Snapshot. This object is passed to
// the snapshots.Create function. For more information about these parameters,
// see the Snapshot object.
type CreateOpts struct {
	// REQUIRED
	VolumeID string
	// OPTIONAL
	Description string
	// OPTIONAL
	Force bool
	// OPTIONAL
	Name string
}

// ToSnapshotCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToSnapshotCreateMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.VolumeID == "" {
		return nil, fmt.Errorf("Required CreateOpts field 'VolumeID' not set.")
	}

	s["volume_id"] = opts.VolumeID

	if opts.Description != "" {
		s["display_description"] = opts.Description
	}
	if opts.Name != "" {
		s["display_name"] = opts.Name
	}
	if opts.Force == true {
		s["force"] = opts.Force
	}

	return map[string]interface{}{"snapshot": s}, nil
}

// Create will create a new Snapshot based on the values in CreateOpts. To
// extract the Snapshot object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	return CreateResult{Common: os.Create(client, opts)}
}

// Delete will delete the existing Snapshot with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) error {
	return os.Delete(client, id)
}

// Get retrieves the Snapshot with the provided ID. To extract the Snapshot
// object from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	return GetResult{Common: os.Get(client, id)}
}

// List returns Snapshots.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return os.List(client, os.ListOpts{})
}
