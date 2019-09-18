package apiversions

import (
	"github.com/gophercloud/gophercloud/pagination"
)

// APIVersion represents an API version for Octavia.
type APIVersion struct {
	// ID is the unique identifier of the API version.
	ID string `json:"id"`

	// Status is the API versions status, e.g. CURRENT, SUPPORTED
	Status string `json:"status"`

	// Updated is the time when the version was added.
	Updated string `json:"updated"`
}

// APIVersionPage is the page returned by a pager when traversing over a
// collection of API versions.
type APIVersionPage struct {
	pagination.SinglePageBase
}

// ExtractAPIVersions takes a collection page, extracts all of the elements,
// and returns them a slice of APIVersion structs
func ExtractAPIVersions(r pagination.Page) ([]APIVersion, error) {
	var s struct {
		Versions []APIVersion `json:"versions"`
	}
	err := (r.(APIVersionPage)).ExtractInto(&s)
	return s.Versions, err
}
