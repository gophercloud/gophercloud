package tenants

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// Tenant is a grouping of users in the identity service.
type Tenant struct {
	// ID is a unique identifier for this tenant.
	ID string `mapstructure:"id"`

	// Name is a friendlier user-facing name for this tenant.
	Name string `mapstructure:"name"`

	// Description is a human-readable explanation of this Tenant's purpose.
	Description string `mapstructure:"description"`

	// Enabled indicates whether or not a tenant is active.
	Enabled bool `mapstructure:"enabled"`
}

// TenantPage is a single page of Tenant results.
type TenantPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Tenants contains any results.
func (page TenantPage) IsEmpty() (bool, error) {
	tenants, err := ExtractTenants(page)
	if err != nil {
		return false, err
	}
	return len(tenants) == 0, nil
}

// NextPageURL extracts the "next" link from the tenants_links section of the result.
func (page TenantPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"tenants_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

// ExtractTenants returns a slice of Tenants contained in a single page of results.
func ExtractTenants(page pagination.Page) ([]Tenant, error) {
	casted := page.(TenantPage).Body
	var response struct {
		Tenants []Tenant `mapstructure:"tenants"`
	}

	fmt.Printf("Decode %#v => %#v\n", casted, response)
	err := mapstructure.Decode(casted, &response)
	return response.Tenants, err
}
