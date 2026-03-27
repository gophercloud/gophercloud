package traits

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// GetResult is the response from a Get operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type GetResult struct {
	gophercloud.ErrResult
}

// TraitsPage contains a single page of all traits from a List call.
type TraitsPage struct {
	pagination.SinglePageBase
}

// IsEmpty satisfies the IsEmpty method of the Page interface. It returns true
// if a List contains no results.
func (r TraitsPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	traits, err := ExtractTraits(r)
	return len(traits) == 0, err
}

// ExtractTraits takes a List result and extracts the collection of traits
// returned by the API.
func ExtractTraits(p pagination.Page) ([]string, error) {
	var s struct {
		Traits []string `json:"traits"`
	}
	err := (p.(TraitsPage)).ExtractInto(&s)
	return s.Traits, err
}
