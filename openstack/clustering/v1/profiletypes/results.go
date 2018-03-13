package profiletypes

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// Extract provides access to the individual node returned by Get and extracts Node
func (r commonResult) Extract() (*ProfileType, error) {
	var s struct {
		ProfileType *ProfileType `json:"profile_type"`
	}
	err := r.ExtractInto(&s)
	return s.ProfileType, err
}

// ExtractProfileTypes provides access to the list of nodes in a page acquired from the ListDetail operation.
func ExtractProfileTypes(r pagination.Page) ([]ProfileType, error) {
	var s struct {
		ProfileTypes []ProfileType `json:"profile_types"`
	}
	err := (r.(ProfileTypePage)).ExtractInto(&s)
	return s.ProfileTypes, err
}

// ProfileTypePage contains a single page of all nodes from a ListDetails call.
type ProfileTypePage struct {
	pagination.LinkedPageBase
}

type ProfileType struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

// IsEmpty determines if a ProfielType contains any results.
func (page ProfileTypePage) IsEmpty() (bool, error) {
	profileTypes, err := ExtractProfileTypes(page)
	return len(profileTypes) == 0, err
}

// OperationPage contains a single page of all profile type operations from a ListOperationDetails call.
type OperationPage struct {
	pagination.LinkedPageBase
}

// Operation represents a detailed operation
type Operation struct {
	Operations map[string]interface{} `json:"operations"`
}

// IsEmpty determines if a Operation contains any results.
func (page OperationPage) IsEmpty() (bool, error) {
	operations, err := ExtractOperations(page)
	return len(operations) == 0, err
}

// ExtractOperations provides access to the list of operations in a page acquired from the ListOperationDetail operation.
func ExtractOperations(r pagination.Page) ([]Operation, error) {
	var s struct {
		Operations []Operation `json:"operations"`
	}
	err := (r.(OperationPage)).ExtractInto(&s)
	if err != nil {
		fmt.Println("Error ExtractOperations")
	}
	return s.Operations, err
}
