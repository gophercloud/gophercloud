package impliedroles

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type ImpliedRole struct {
	PriorRole struct {
		ID string `json:"id"`

		Name string `json:"name"`

		Links map[string]interface{} `json:"links"`
	} `json:"prior_role"`

	Implies []struct {
		ID string `json:"id"`

		Name string `json:"name"`

		Links map[string]interface{} `json:"links"`
	} `json:"implies"`
}

type CreateImpliedRole struct {
	PriorRole struct {
		ID string `json:"id"`

		Name string `json:"name"`

		Links map[string]interface{} `json:"links"`
	} `json:"prior_role"`

	Implies struct {
		ID string `json:"id"`

		Name string `json:"name"`

		Links map[string]interface{} `json:"links"`
	} `json:"implies"`
}

type createImpliedRoleResult struct {
	gophercloud.Result
}

type impliedRoleResult struct {
	gophercloud.Result
}

// CreateResult is the response from a Create operation. Call its ExtractErr to
// determine if the request succeded or failed.
type CreateResult struct {
	gophercloud.Result
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ImpliedRolePage is a single page of Role results.
type ImpliedRolePage struct {
	pagination.LinkedPageBase
}

// CreateImpliedRolePage
type CreateImpliedRolePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of ImpliedRoles contains any results.
func (r ImpliedRolePage) IsEmpty() (bool, error) {
	Impliedroles, err := ExtractImpliedRoles(r)
	return len(Impliedroles) == 0, err
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

// ExtractProjects returns a slice of ImpliedRoles contained in a single page of
// results.
func ExtractImpliedRoles(r pagination.Page) ([]ImpliedRole, error) {
	var s struct {
		ImpliedRoles []ImpliedRole `json:"role_inferences"`
	}
	err := (r.(ImpliedRolePage)).ExtractInto(&s)
	return s.ImpliedRoles, err
}

// Extract interprets any impliedRoleResults as a Role.
func (r impliedRoleResult) Extract() (*ImpliedRole, error) {
	var s struct {
		ImpliedRole *ImpliedRole `json:"role_inferences"`
	}
	err := r.ExtractInto(&s)
	return s.ImpliedRole, err
}

// Extract interprets any createImpliedRoleResult as a role
func (r createImpliedRoleResult) Extract() (*CreateImpliedRole, error) {
	var s struct {
		CreateImpliedRole *CreateImpliedRole `json:"role_inference"`
	}
	err := r.ExtractInto(&s)
	return s.CreateImpliedRole, err
}
