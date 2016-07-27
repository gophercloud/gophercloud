package baymodels

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a baymodel resource.
func (r commonResult) Extract() (*BayModel, error) {
	var s *BayModel
	err := r.ExtractInto(&s)
	return s, err
}

// Represents a template for a Bay
type BayModel struct {
	// UUID for the baymodel
	ID string `json:"uuid"`

	// Human-readable name for the baymodel. Might not be unique.
	Name string `json:"name"`

	// The type of container orchestration engine used by the bay.
	COE string `json:"coe"`

	// The flavor used by nodes in the bay.
	FlavorID string `json:"flavor_id"`

	// The image used by nodes in the bay.
	ImageID string `json:"image_id"`

	// Specifies if the bay should use TLS certificates.
	TLSDisabled bool `json:"tls_disabled"`

	// The KeyPair used by the bay.
	KeyPairID string `json:"keypair_id"`
}

// BayModelPage is the page returned by a pager when traversing over a
// collection of baymodels.
type BayModelPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of baymodels has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BayModelPage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Next, nil
}

// IsEmpty checks whether a BayModelPage struct is empty.
func (r BayModelPage) IsEmpty() (bool, error) {
	is, err := ExtractBayModels(r)
	return len(is) == 0, err
}

// ExtractBayModels accepts a Page struct, specifically a BayModelPage struct,
// and extracts the elements into a slice of BayModel structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBayModels(r pagination.Page) ([]BayModel, error) {
	var s struct {
		BayModels []BayModel `json:"baymodels"`
	}
	err := (r.(BayModelPage)).ExtractInto(&s)
	return s.BayModels, err
}
