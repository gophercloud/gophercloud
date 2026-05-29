package osinherit

import "github.com/gophercloud/gophercloud/v2"

// AssignmentResult represents the result of an assign operation.
// Call ExtractErr method to determine if the request succeeded or failed.
type AssignmentResult struct {
	gophercloud.ErrResult
}

// ValidateResult represents the result of an validate operation.
// Call ExtractErr method to determine if the request succeeded or failed.
type ValidateResult struct {
	gophercloud.ErrResult
}

// UnassignmentResult represents the result of an unassign operation.
// Call ExtractErr method to determine if the request succeeded or failed.
type UnassignmentResult struct {
	gophercloud.ErrResult
}
