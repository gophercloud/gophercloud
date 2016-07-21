package bays

import (
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a bay resource.
func (r commonResult) Extract() (*Bay, error) {
	var bay *Bay
	err := r.ExtractInto(&bay)
	return bay, err
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// Represents a Container Orchestration Engine Bay, i.e. a cluster
type Bay struct {
	// UUID for the bay
	ID string `json:"uuid"`

	// Human-readable name for the bay. Might not be unique.
	Name string `json:"name"`

	// Indicates whether bay is currently operational. Possible values include:
	// CREATE_IN_PROGRESS, CREATE_FAILED, CREATE_COMPLETE, UPDATE_IN_PROGRESS, UPDATE_FAILED, UPDATE_COMPLETE,
	// DELETE_IN_PROGRESS, DELETE_FAILED, DELETE_COMPLETE, RESUME_COMPLETE, RESTORE_COMPLETE, ROLLBACK_COMPLETE,
	// SNAPSHOT_COMPLETE, CHECK_COMPLETE, ADOPT_COMPLETE.
	Status string `json:"status"`

	// The number of nodes in the bay.
	Nodes int `json:"node_count"`

	// The UUID of the baymodel used to generate the bay.
	BayModelID string `json:"baymodel_id"`
}
