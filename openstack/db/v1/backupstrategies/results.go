package backupstrategies

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// BackupStrategy represents a database backup strategy.
type BackupStrategy struct {
	ProjectID      string `json:"project_id"`
	InstanceID     string `json:"instance_id"`
	Backend        string
	SwiftContainer string `json:"swift_container"`
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	gophercloud.Result
}

// Extract retrieves a BackupStrategy resource from an operation result.
func (r CreateResult) Extract() (*BackupStrategy, error) {
	var s struct {
		BackupStrategy *BackupStrategy `json:"backup_strategy"`
	}
	err := r.ExtractInto(&s)
	return s.BackupStrategy, err
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// BackupStrategyPage represents a page of backup strategies.
type BackupStrategyPage struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a BackupStrategyPage is empty.
func (r BackupStrategyPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	strategies, err := ExtractBackupStrategies(r)
	return len(strategies) == 0, err
}

// ExtractBackupStrategies retrieves a slice of BackupStrategy structs from a
// paginated collection.
func ExtractBackupStrategies(r pagination.Page) ([]BackupStrategy, error) {
	var s struct {
		BackupStrategies []BackupStrategy `json:"backup_strategies"`
	}
	err := (r.(BackupStrategyPage)).ExtractInto(&s)
	return s.BackupStrategies, err
}
