package flavorprofiles

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// FlavorProfile provide metadata such as provider, toplogy and instance flavor.
type FlavorProfile struct {
	// The unique ID for the Flavor
	ID string `json:"id"`

	// Human-readable name for the Flavor. Does not have to be unique.
	Name string `json:"name"`

	// Name of the provider
	ProviderName string `json:"provider_name"`

	// Flavor data
	FlavorData string `json:"flavor_data"`
}

// FlavorProfilePage is the page returned by a pager when traversing over a
// collection of flavor profiles.
type FlavorProfilePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of flavor profiles has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r FlavorProfilePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"flavorprofiles_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FlavorProfilePage struct is empty.
func (r FlavorProfilePage) IsEmpty() (bool, error) {
	is, err := ExtractFlavorProfiles(r)
	return len(is) == 0, err
}

// ExtractFlavorProfiles accepts a Page struct, specifically a FlavorProfilePage
// struct, and extracts the elements into a slice of FlavorProfile structs. In
// other words, a generic collection is mapped into a relevant slice.
func ExtractFlavorProfiles(r pagination.Page) ([]FlavorProfile, error) {
	var s struct {
		FlavorProfiles []FlavorProfile `json:"flavorprofiles"`
	}
	err := (r.(FlavorProfilePage)).ExtractInto(&s)
	return s.FlavorProfiles, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a flavor profile.
func (r commonResult) Extract() (*FlavorProfile, error) {
	var s struct {
		FlavorProfile *FlavorProfile `json:"flavorprofile"`
	}
	err := r.ExtractInto(&s)
	return s.FlavorProfile, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a FlavorProfile.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a FlavorProfile.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a FlavorProfile.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
