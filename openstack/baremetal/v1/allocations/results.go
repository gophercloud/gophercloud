package allocations

import (
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Allocation struct {
	// The UUID for the resource.
	UUID string `json:"uuid"`

	// A list of UUIDs of the nodes that are candidates for this allocation.
	CandidateNodes []string `json:"candidate_nodes"`

	// The error message for the allocation if it is in the error state, null otherwise.
	LastError string `json:"last_error"`

	// The unique name of the allocation.
	Name string `json:"name"`

	// The UUID of the node assigned to the allocation. Will be null if a node is not yet assigned.
	NodeUUID string `json:"node_uuid"`

	// The current state of the allocation. One of: allocation, active, error
	State string `json:"state"`

	// The resource class requested for the allocation.
	ResourceClass string `json:"resource_class"`

	// The list of the traits requested for the allocation.
	Traits []string `json:"traits"`

	// A set of one or more arbitrary metadata key and value pairs.
	Extra map[string]string `json:"extra"`

	// The UTC date and time when the resource was created, ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`

	// The UTC date and time when the resource was updated, ISO 8601 format. May be “null”.
	UpdatedAt time.Time `json:"updated_at"`

	// A list of relative links. Includes the self and bookmark links.
	Links []any `json:"links"`
}

type allocationResult struct {
	gophercloud.Result
}

func (r allocationResult) Extract() (*Allocation, error) {
	var s Allocation
	err := r.ExtractInto(&s)
	return &s, err
}

func (r allocationResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "")
}

func ExtractAllocationsInto(r pagination.Page, v any) error {
	return r.(AllocationPage).ExtractIntoSlicePtr(v, "allocations")
}

// AllocationPage abstracts the raw results of making a List() request against
// the API.
type AllocationPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Allocation results.
func (r AllocationPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	s, err := ExtractAllocations(r)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r AllocationPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"allocations_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractAllocations interprets the results of a single page from a List() call,
// producing a slice of Allocation entities.
func ExtractAllocations(r pagination.Page) ([]Allocation, error) {
	var s []Allocation
	err := ExtractAllocationsInto(r, &s)
	return s, err
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a Allocation.
type GetResult struct {
	allocationResult
}

// CreateResult is the response from a Create operation.
type CreateResult struct {
	allocationResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
