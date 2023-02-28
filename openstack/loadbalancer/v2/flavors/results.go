package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type Flavor struct {
	// The unique ID for the Flavor
	ID string `json:"id"`

	// Human-readable name for the Flavor. Does not have to be unique.
	Name string `json:"name"`

	// Human-readable description for the Flavor.
	Description string `json:"description"`

	// Status of the Flavor.
	Enabled bool `json:"enabled"`

	// Flavor Profile apply to this Flavor.
	FlavorProfileId string `json:"flavor_profile_id"`
}

type FlavorPage struct {
	pagination.LinkedPageBase
}

func (r FlavorPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"flavors_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

func (r FlavorPage) IsEmpty() (bool, error) {
	is, err := ExtractFlavors(r)
	return len(is) == 0, err
}

func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorPage)).ExtractInto(&s)
	return s.Flavors, err
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Flavor, error) {
	var s struct {
		Flavor *Flavor `json:"flavor"`
	}
	err := r.ExtractInto(&s)
	return s.Flavor, err
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
