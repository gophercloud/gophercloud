package receivers

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToReceiverCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name      string                 `json:"name" required:"true"`
	ClusterID string                 `json:"cluster_id,omitempty"`
	Type      string                 `json:"type" required:"true"`
	Action    string                 `json:"action,omitempty"`
	Actor     map[string]interface{} `json:"actor,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

// ToReceiverCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToReceiverCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "receiver")
}

// Create requests the creation of a new receiver.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToReceiverCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return
}

// Get retrieves details of a single receiver. Use Extract to convert its result into a Receiver.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
