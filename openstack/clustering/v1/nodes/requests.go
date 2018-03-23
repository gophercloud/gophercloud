package nodes

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// UpdateOpts params
type UpdateOpts struct {
	Node      map[string]interface{} `json:"-"`
	Name      string                 `json:"name,omitempty"`
	ProfileID string                 `json:"profile_id,omitempty"`
	Role      string                 `json:"role,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateOptsBuilder params
type UpdateOptsBuilder interface {
	ToNodeUpdateMap() (map[string]interface{}, error)
}

// ToClusterUpdateMap constructs a request body from CreateOpts.
func (opts UpdateOpts) ToNodeUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "node")
}

// Update requests the update of a node.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNodeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
