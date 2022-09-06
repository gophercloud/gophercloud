package limits

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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

// A limit is the limit that override the registered limit for each project.
type Limit struct {
	// ID is the unique ID of the limit.
	ID string `json:"id"`

	// RegionID is the ID of the region where the limit is applied.
	RegionID string `json:"region_id"`

	// ProjectID is the ID of the project where the limit is applied.
	ProjectID string `json:"project_id"`

	// DomainID is the ID of the domain where the limit is applied.
	DomainID string `json:"domain_id"`

	// ServiceID is the ID of the service where the limit is applied.
	ServiceID string `json:"service_id"`

	// Description of the limit.
	Description string `json:"description"`

	// ResourceName is the name of the resource that the limit is applied to.
	ResourceName string `json:"resource_name"`

	// ResourceLimit is the override limit.
	ResourceLimit int `json:"resource_limit"`

	// Links contains referencing links to the limit.
	Links map[string]interface{} `json:"links"`
}

// LimitPage is a single page of Limit results.
type LimitPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Limits contains any results.
func (r LimitPage) IsEmpty() (bool, error) {
	limits, err := ExtractLimits(r)
	return len(limits) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r LimitPage) NextPageURL() (string, error) {
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

// ExtractLimits returns a slice of Limits contained in a single page of
// results.
func ExtractLimits(r pagination.Page) ([]Limit, error) {
	var s struct {
		Limits []Limit `json:"limits"`
	}
	err := (r.(LimitPage)).ExtractInto(&s)
	return s.Limits, err
}
