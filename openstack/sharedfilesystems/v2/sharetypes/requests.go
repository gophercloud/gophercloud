package sharetypes

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareTypeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a ShareType. This object is
// passed to the sharetypes.Create function. For more information about
// these parameters, see the ShareType object.
type CreateOpts struct {
	// The share type name
	Name string `json:"name" required:"true"`
	// Indicates whether a share type is publicly accessible
	IsPublic bool `json:"os-share-type-access:is_public"`
	// The extra specifications for the share type
	ExtraSpecs ExtraSpecsOpts `json:"extra_specs" required:"true"`
}

// ExtraSpecsOpts represent the extra specifications that can be selected for a share type
type ExtraSpecsOpts struct {
	// An extra specification that defines the driver mode for share server, or storage, life cycle management
	DriverHandlesShareServers bool `json:"driver_handles_share_servers" required:"true"`
	// An extra specification that filters back ends by whether they do or do not support share snapshots
	SnapshotSupport *bool `json:"snapshot_support,omitempty"`
}

// ToShareTypeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareTypeCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "share_type")
}

// Create will create a new ShareType based on the values in CreateOpts. To
// extract the ShareType object from the response, call the Extract method
// on the CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareTypeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will delete the existing ShareType with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
