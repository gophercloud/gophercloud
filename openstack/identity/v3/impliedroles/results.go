package impliedroles

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type ImpliedRole struct {
	PriorRole struct {
		ID string `json:"id"`

		Name string `json:"name"`

		// // Extra is a collection of miscellaneous key/values.
		// Extra map[string]interface{} `json:"-"`
	}

	Implies struct {
		ID string `json:"id"`

		Name string `json:"name"`

		// Extra is a collection of miscellaneous key/values.
		// Extra map[string]interface{} `json:"-"`
	}

	// Extra is a collection of miscellaneous key/values.
	Extra map[string]interface{} `json:"-"`
}

func (r *ImpliedRole) UnmarshalJSON(b []byte) error {
	type tmp ImpliedRole
	var s struct {
		tmp
		Extra map[string]interface{} `json:"extra"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ImpliedRole(s.tmp)

	// Collect other fields and bundle them into Extra
	// but only if a field titled "extra" wasn't sent.
	if s.Extra != nil {
		r.Extra = s.Extra
	} else {
		var result interface{}
		err := json.Unmarshal(b, &result)
		if err != nil {
			return err
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			r.Extra = gophercloud.RemainingKeys(ImpliedRole{}, resultMap)
		}
	}

	return err
}

type impliedRoleResult struct {
	gophercloud.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as a Role.
type GetResult struct {
	impliedRoleResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a Role
type CreateResult struct {
	impliedRoleResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// RolePage is a single page of Role results.
type ImpliedRolePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Roles contains any results.
func (r ImpliedRolePage) IsEmpty() (bool, error) {
	roles, err := ExtractImpliedRoles(r)
	return len(roles) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r ImpliedRolePage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractProjects returns a slice of Roles contained in a single page of
// results.
func ExtractImpliedRoles(r pagination.Page) ([]ImpliedRole, error) {
	var s struct {
		ImpliedRoles []ImpliedRole `json:"role_inference"`
	}
	err := (r.(ImpliedRolePage)).ExtractInto(&s)
	return s.ImpliedRoles, err
}

// Extract interprets any roleResults as a Role.
func (r impliedRoleResult) Extract() (*ImpliedRole, error) {
	var s struct {
		Role *ImpliedRole `json:"role"`
	}
	err := r.ExtractInto(&s)
	return s.Role, err
}
