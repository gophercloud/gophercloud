package registeredlimits

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// A model describing the configured enforcement model used by the deployment.
type EnforcementModel struct {
	// The name of the enforcement model.
	Name string `json:"name"`

	// A short description of the enforcement model used.
	Description string `json:"description"`
}

// EnforcementModelResult is the response from a GetEnforcementModel operation. Call its Extract method
// to interpret it as a EnforcementModel.
type EnforcementModelResult struct {
	gophercloud.Result
}

// Extract interprets EnforcementModelResult as a EnforcementModel.
func (r EnforcementModelResult) Extract() (*EnforcementModel, error) {
	var out struct {
		Model *EnforcementModel `json:"model"`
	}
	err := r.ExtractInto(&out)
	return out.Model, err
}

// A registered limit is the limit that is default for all projects.
type RegisteredLimit struct {
	// ID is the unique ID of the limit.
	ID string `json:"id"`

	// RegionID is the ID of the region where the limit is applied.
	RegionID string `json:"region_id"`

	// ServiceID is the ID of the service where the limit is applied.
	ServiceID string `json:"service_id"`

	// Description of the limit.
	Description string `json:"description"`

	// ResourceName is the name of the resource that the limit is applied to.
	ResourceName string `json:"resource_name"`

	// DefaultLimit is the default limit.
	DefaultLimit int `json:"default_limit"`

	// Links contains referencing links to the limit.
	Links map[string]any `json:"links"`
}

// A LimitsOutput is an array of limits returned by List and BatchCreate operations
type RegisteredLimitsOutput struct {
	RegisteredLimits []RegisteredLimit `json:"registered_limits"`
}

// A RegisteredLimitOutput is an encapsulated Limit returned by Get and Update operations
type RegisteredLimitOutput struct {
	RegisteredLimit *RegisteredLimit `json:"registered_limit"`
}

// RegisteredLimitPage is a single page of Registered Limit results.
type RegisteredLimitPage struct {
	pagination.LinkedPageBase
}

type commonResult struct {
	gophercloud.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as a RegisteredLimit.
type GetResult struct {
	commonResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a Registered Limits.
type CreateResult struct {
	gophercloud.Result
}

// UpdateResult is the result of an Update request. Call its Extract method to
// interpret it as a Limit.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a Delete request. Call its ExtractErr method to
// determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// IsEmpty determines whether or not a page of Limits contains any results.
func (r RegisteredLimitPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	registered_limits, err := ExtractRegisteredLimits(r)
	return len(registered_limits) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r RegisteredLimitPage) NextPageURL() (string, error) {
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

// ExtractRegisteredLimits returns a slice of Registered Limits contained in a single page of
// results.
func ExtractRegisteredLimits(r pagination.Page) ([]RegisteredLimit, error) {
	var out RegisteredLimitsOutput
	err := (r.(RegisteredLimitPage)).ExtractInto(&out)
	return out.RegisteredLimits, err
}

// Extract interprets CreateResult as slice of RegisteredLimits.
func (r CreateResult) Extract() ([]RegisteredLimit, error) {
	var out RegisteredLimitsOutput
	err := r.ExtractInto(&out)
	return out.RegisteredLimits, err
}

// Extract interprets any commonResult as a RegisteredLimit.
func (r commonResult) Extract() (*RegisteredLimit, error) {
	var out RegisteredLimitOutput
	err := r.ExtractInto(&out)
	return out.RegisteredLimit, err
}
