package quotas

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Quota resource.
func (r commonResult) Extract() (*Quota, error) {
	var s struct {
		Quota *Quota `json:"quotas"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetResult struct {
	commonResult
}

// Quota contains key manager quotas for a project.
type Quota struct {
	// Secrets represents the number of secrets. A "-1" value means no limit.
	Secrets *int `json:"secrets"`

	// Orders represents the number of orders. A "-1" value means no limit.
	Orders *int `json:"orders"`

	// Containers represents the number of containers. A "-1" value means no limit.
	Containers *int `json:"containers"`

	// Consumers represents the number of consumers. A "-1" value means no limit.
	Consumers *int `json:"consumers"`

	// CAS represents the number of cas. A "-1" value means no limit.
	CAS *int `json:"cas"`
}

type commonProjectResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Quota resource.
func (r commonProjectResult) Extract() (*Quota, error) {
	var s struct {
		Quota *Quota `json:"project_quotas"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// GetProjectResult represents the result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetProjectResult struct {
	commonProjectResult
}

type ProjectQuota struct {
	ProjectID string `json:"project_id"`
	Quota     Quota  `json:"project_quotas"`
}

// ProjectQuotaPage is a single page of quotas results.
type ProjectQuotaPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of quotas contains any results.
func (r ProjectQuotaPage) IsEmpty() (bool, error) {
	quotas, err := ExtractQuotas(r)
	return len(quotas) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r ProjectQuotaPage) NextPageURL() (string, error) {
	var s struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Next, err
}

// ExtractQuotas returns a slice of ProjectQuota contained in a single page of
// results.
func ExtractQuotas(r pagination.Page) ([]ProjectQuota, error) {
	var s struct {
		Quotas []ProjectQuota `json:"project_quotas"`
	}
	err := (r.(ProjectQuotaPage)).ExtractInto(&s)
	return s.Quotas, err
}
