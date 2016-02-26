package flavors

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ErrCannotInterpret is returned by an Extract call if the response body doesn't have the expected structure.
var ErrCannotInterpet = errors.New("Unable to interpret a response body.")

// GetResult temporarily holds the response from a Get call.
type GetResult struct {
	gophercloud.Result
}

// Extract provides access to the individual Flavor returned by the Get function.
func (r GetResult) Extract() (*Flavor, error) {
	var s struct {
		Flavor *Flavor `json:"flavor"`
	}
	err := r.ExtractInto(&s)
	return s.Flavor, err
}

// Flavor records represent (virtual) hardware configurations for server resources in a region.
type Flavor struct {
	// The Id field contains the flavor's unique identifier.
	// For example, this identifier will be useful when specifying which hardware configuration to use for a new server instance.
	ID string `json:"id"`

	// The Disk and RA< fields provide a measure of storage space offered by the flavor, in GB and MB, respectively.
	Disk int `json:"disk"`
	RAM  int `json:"ram"`

	// The Name field provides a human-readable moniker for the flavor.
	Name string `json:"name"`

	RxTxFactor float64 `json:"rxtx_factor"`

	// Swap indicates how much space is reserved for swap.
	// If not provided, this field will be set to 0.
	Swap int `json:"swap"`

	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs int `json:"vcpus"`
}

// FlavorPage contains a single page of the response from a List call.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a page contains any results.
func (page FlavorPage) IsEmpty() (bool, error) {
	flavors, err := ExtractFlavors(page)
	return len(flavors) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page FlavorPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"flavors_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractFlavors provides access to the list of flavors in a page acquired from the List operation.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorPage)).ExtractInto(&s)
	return s.Flavors, err
}
