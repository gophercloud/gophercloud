package rbacpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// PolicyAction maps to Action for the RBAC policy.
// Which allows access_as_external or access_as_shared.
type PolicyAction string

const (
	// ActionAccessExternal returns Action for the RBAC policy as access_as_external.
	ActionAccessExternal PolicyAction = "access_as_external"

	// ActionAccessShared returns Action for the RBAC policy as access_as_shared.
	ActionAccessShared PolicyAction = "access_as_shared"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRBACPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a rbac-policy.
type CreateOpts struct {
	Action       PolicyAction `json:"action" required:"true"`
	ObjectType   string       `json:"object_type" required:"true"`
	TargetTenant string       `json:"target_tenant" required:"true"`
	ObjectID     string       `json:"object_id" required:"true"`
}

// ToRBACPolicyCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRBACPolicyCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "rbac_policy")
}

// Create accepts a CreateOpts struct and creates a new rbac-policy using the values
// provided.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// rbac-policy.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRBACPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}
