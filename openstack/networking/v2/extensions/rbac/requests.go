package rbac

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRbacCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a rbac_policy.
type CreateOpts struct {
	Action       string `json:"action"`
	ObjectType   string `json:"object_type"`
	TargetTenant string `json:"target_tenant"`
	ObjectID     string `json:"object_id"`
}

// ToRbacCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRbacCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "rbac_policy")
}

// Create accepts a CreateOpts struct and creates a new rbac-policy using the values
// provided.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// rbac-policy.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRbacCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}
