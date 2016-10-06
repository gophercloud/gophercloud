package shares

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts used for ...
type CreateOpts struct {
	ShareProto         string            `json:"share_proto"`
	Size               int               `json:"size"`
	Name               string            `json:"name,omitempty"`
	Description        string            `json:"description,omitempty"`
	DisplayName        string            `json:"display_name,omitempty"`
	DisplayDescription string            `json:"display_description,omitempty"`
	ShareType          string            `json:"share_type,omitempty"`
	VolumeType         string            `json:"volume_type,omitempty"`
	SnapshotID         string            `json:"snapshot_id,omitempty"`
	IsPublic           bool              `json:"is_public,omitempty"`
	MetaData           map[string]string `json:"metadata,omitempty"`
	ShareNetworkID     string            `json:"share_network_id,omitempty"`
	ConsistencyGroupID string            `json:"consistency_group_id,omitempty"`
	AvailabilityZone   string            `json:"availability_zone,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "share")
}

// Create will create a new Share based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}
