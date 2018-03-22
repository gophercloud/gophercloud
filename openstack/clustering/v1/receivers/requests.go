package receivers

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// UpdateOpts params
type UpdateOpts struct {
	Name   string                 `json:"name,omitempty"`
	Action string                 `json:"action,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type UpdateOptsBuilder interface {
	ToReceiverUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToReceiverUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "receiver")
}

// Update requests the update of a receiver.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToReceiverUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
