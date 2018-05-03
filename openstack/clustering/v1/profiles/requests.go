package profiles

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder for options used for creating a profile.
type CreateOptsBuilder interface {
	ToProfileCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used for creating a profile
type CreateOpts struct {
	Name     string                 `json:"name" required:"true"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Spec     Spec                   `json:"spec" required:"true"`
}

// ToProfileCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToProfileCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "profile")
}

// Create requests the creation of a new profile on the server.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})

	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
