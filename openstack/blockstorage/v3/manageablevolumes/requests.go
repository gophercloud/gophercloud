package manageablevolumes

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// ManageExistingOptsBuilder allows extentions to add additional parameters to the ManageExisting request.
type ManageExistingOptsBuilder interface {
	ToManageExistingMap() (map[string]any, error)
}

// ManageExistingOpts contains options for managing a existing volume.
// This object is passed to the volumes.ManageExisting function.
// For more information about the parameters, see the Volume object and OpenStack BlockStorage API Guide.
type ManageExistingOpts struct {
	// The OpenStack Block Storage host where the existing resource resides.
	// Optional only if cluster field is provided.
	Host string `json:"host,omitempty"`
	// The OpenStack Block Storage cluster where the resource resides.
	// Optional only if host field is provided.
	Cluster string `json:"cluster,omitempty"`
	// A reference to the existing volume.
	// The internal structure of this reference depends on the volume driver implementation.
	// For details about the required elements in the structure, see the documentation for the volume driver.
	Ref map[string]string `json:"ref,omitempty"`
	// Human-readable display name for the volume.
	Name string `json:"name,omitempty"`
	// The availability zone.
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// Human-readable description for the volume.
	Description string `json:"description,omitempty"`
	// The associated volume type
	VolumeType string `json:"volume_type,omitempty"`
	// Indicates whether this is a bootable volume.
	Bootable bool `json:"bootable,omitempty"`
	// One or more metadata key and value pairs to associate with the volume.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ToManageExistingMap assembles a request body based on the contents of a ManageExistingOpts.
func (opts ManageExistingOpts) ToManageExistingMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "volume")
}

// ManageExisting will manage an existing volume based on the values in ManageExistingOpts.
// To extract the Volume object from response, call the Extract method on the ManageExistingResult.
func ManageExisting(ctx context.Context, client *gophercloud.ServiceClient, opts ManageExistingOptsBuilder) (r ManageExistingResult) {
	b, err := opts.ToManageExistingMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
