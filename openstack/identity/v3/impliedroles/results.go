package impliedroles

import (
	"github.com/gophercloud/gophercloud"
)

type Role struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Links map[string]interface{} `json:"links"`
}

type ImpliedRole struct {
	PriorRole Role `json:"prior_role"`

	Implies []Role `json:"implies"`
}

type GetImpliedRole struct {
	RoleInference ImpliedRole `json:"role_inference"`

	Links map[string]interface{} `json:"links"`
}

type GetImpliedRoleResult struct {
	gophercloud.Result
}

// Extract interprets any impliedRoleResults as a Role.
func (r GetImpliedRoleResult) Extract() (*GetImpliedRole, error) {
	var s GetImpliedRole
	err := r.ExtractInto(&s)
	return &s, err
}
