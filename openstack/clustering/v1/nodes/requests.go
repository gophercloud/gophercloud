package nodes

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToNodeCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	Role      string                 `json:"role,omitempty"`
	ProfileID string                 `json:"profile_id" required:"true"`
	ClusterID string                 `json:"cluster_id,omitempty"`
	Name      string                 `json:"name" required:"true"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ToNodeCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToNodeCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "node")
}

// Create requests the creation of a new node.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNodeCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
