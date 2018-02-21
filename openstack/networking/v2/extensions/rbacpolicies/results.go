package rbacpolicies

import (
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts RBAC Policy resource.
func (r commonResult) Extract() (*RBACPolicy, error) {
	var s RBACPolicy
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "rbac_policy")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a RBAC Policy.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a RBAC Policy.
type GetResult struct {
	commonResult
}

// RBACPolicy represents a RBAC policy.
type RBACPolicy struct {
	// UUID of the RBAC policy.
	ID string `json:"id"`

	// Action for the RBAC policy which is access_as_external or access_as_shared.
	Action PolicyAction `json:"action"`

	// ObjectID is the ID of the object_type resource.
	// An object_type of network returns a network ID and
	// object_type of qos-policy returns a QoS ID.
	ObjectID string `json:"object_id"`

	// ObjectType is the type of the object that the RBAC policy affects.
	// Types include qos-policy or network.
	ObjectType string `json:"object_type"`

	// TenantID is the ID of the project that owns the resource.
	TenantID string `json:"tenant_id"`

	// TargetTenant is the ID of the tenant to which the RBAC policy will be enforced.
	TargetTenant string `json:"target_tenant"`

	// ProjectID is the ID of the project.
	ProjectID string `json:"project_id"`
}
