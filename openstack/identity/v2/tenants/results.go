package tenants

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Tenant is a grouping of users in the identity service.
type Tenant struct {
	// ID is a unique identifier for this tenant.
	ID string `json:"id"`

	// Name is a friendlier user-facing name for this tenant.
	Name string `json:"name"`

	// Description is a human-readable explanation of this Tenant's purpose.
	Description string `json:"description"`

	// Enabled indicates whether or not a tenant is active.
	Enabled bool `json:"enabled"`
}

// TenantPage is a single page of Tenant results.
type TenantPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Tenants contains any results.
func (page TenantPage) IsEmpty() (bool, error) {
	tenants, err := ExtractTenants(page)
	return len(tenants) == 0, err
}

// NextPageURL extracts the "next" link from the tenants_links section of the result.
func (page TenantPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"tenants_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractTenants returns a slice of Tenants contained in a single page of results.
func ExtractTenants(page pagination.Page) ([]Tenant, error) {
	r := page.(TenantPage)
	var s struct {
		Tenants []Tenant `json:"tenants"`
	}
	err := r.ExtractInto(&s)
	return s.Tenants, err
}
