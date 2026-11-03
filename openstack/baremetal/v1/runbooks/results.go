package runbooks

import (
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type runbookResult struct {
	gophercloud.Result
}

// Extract interprets any runbookResult as a Runbook, if possible.
func (r runbookResult) Extract() (*Runbook, error) {
	var s Runbook
	err := r.ExtractInto(&s)
	return &s, err
}

func (r runbookResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "")
}

// ExtractRunbooksInto extracts the runbooks of a page into v.
func ExtractRunbooksInto(r pagination.Page, v any) error {
	return r.(RunbookPage).ExtractIntoSlicePtr(v, "runbooks")
}

// Runbook represents a runbook in the OpenStack Bare Metal API.
type Runbook struct {
	// UUID for the resource.
	UUID string `json:"uuid"`

	// Name is the unique name of the runbook.
	Name string `json:"name"`

	// Steps are the steps the runbook runs.
	Steps []Step `json:"steps"`

	// DisableRamdisk shows whether the steps run without the IPA ramdisk.
	DisableRamdisk bool `json:"disable_ramdisk"`

	// Extra is a set of arbitrary metadata key and value pairs.
	Extra map[string]any `json:"extra"`

	// Public shows whether every project can use the runbook.
	Public bool `json:"public"`

	// Owner is the project that owns the runbook. It is empty for
	// public runbooks.
	Owner string `json:"owner"`

	// Description is a free-text description of the runbook.
	// This requires the client to be set to microversion 1.112 or later.
	Description *string `json:"description"`

	// Traits are the node traits that allow a node to run the runbook.
	// This requires the client to be set to microversion 1.112 or later.
	Traits []string `json:"traits"`

	// The UTC date and time when the resource was created, ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`

	// The UTC date and time when the resource was updated, ISO 8601 format.
	// May be “null”.
	UpdatedAt time.Time `json:"updated_at"`

	// A list of relative links. Includes the self and bookmark links.
	Links []any `json:"links"`
}

// RunbookPage abstracts the raw results of making a List() request against
// the API.
type RunbookPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Runbook results.
func (r RunbookPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	s, err := ExtractRunbooks(r)
	return len(s) == 0, err
}

// NextPageURL returns the URL of the next page, which Ironic returns as
// "next" when more results are available.
func (r RunbookPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	return s.Next, err
}

// ExtractRunbooks interprets the results of a single page from a List() call,
// producing a slice of Runbook entities.
func ExtractRunbooks(r pagination.Page) ([]Runbook, error) {
	var s []Runbook
	err := ExtractRunbooksInto(r, &s)
	return s, err
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a Runbook.
type GetResult struct {
	runbookResult
}

// CreateResult is the response from a Create operation. Call its Extract
// method to interpret it as a Runbook.
type CreateResult struct {
	runbookResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Runbook.
type UpdateResult struct {
	runbookResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ListTraitsResult is the response from a ListTraits operation. Call its
// Extract method to interpret it as a list of traits.
type ListTraitsResult struct {
	gophercloud.Result
}

// Extract interprets a ListTraitsResult as a list of traits.
func (r ListTraitsResult) Extract() ([]string, error) {
	var s struct {
		Traits []string `json:"traits"`
	}
	err := r.ExtractInto(&s)
	return s.Traits, err
}

// SetTraitsResult is the response from a SetTraits operation. Call its
// ExtractErr method to determine if the call succeeded or failed.
type SetTraitsResult struct {
	gophercloud.ErrResult
}

// AddTraitResult is the response from an AddTrait operation. Call its
// ExtractErr method to determine if the call succeeded or failed.
type AddTraitResult struct {
	gophercloud.ErrResult
}

// DeleteTraitsResult is the response from a DeleteTraits operation. Call its
// ExtractErr method to determine if the call succeeded or failed.
type DeleteTraitsResult struct {
	gophercloud.ErrResult
}

// DeleteTraitResult is the response from a DeleteTrait operation. Call its
// ExtractErr method to determine if the call succeeded or failed.
type DeleteTraitResult struct {
	gophercloud.ErrResult
}
