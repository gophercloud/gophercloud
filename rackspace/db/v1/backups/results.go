package backups

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/datastores"
)

// Status represents the various states a Backup can be in.
type Status string

// Enum types for the status.
const (
	StatusNew          Status = "NEW"
	StatusBuilding     Status = "BUILDING"
	StatusCompleted    Status = "COMPLETED"
	StatusFailed       Status = "FAILED"
	StatusDeleteFailed Status = "DELETE_FAILED"
)

// Backup represents a Backup API resource.
type Backup struct {
	Description string
	ID          string
	InstanceID  string `json:"instance_id" mapstructure:"instance_id"`
	LocationRef string
	Name        string
	ParentID    string `json:"parent_id" mapstructure:"parent_id"`
	Size        float64
	Status      Status
	Created     string
	Updated     string
	Datastore   datastores.DatastorePartial
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

type commonResult struct {
	gophercloud.Result
}

// Extract will retrieve a Backup struct from an operation's result.
func (r commonResult) Extract() (*Backup, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Backup Backup `mapstructure:"backup"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.Backup, err
}

// BackupPage represents a page of backups.
type BackupPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether an BackupPage struct is empty.
func (r BackupPage) IsEmpty() (bool, error) {
	is, err := ExtractBackups(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractBackups will retrieve a slice of Backup structs from a paginated collection.
func ExtractBackups(page pagination.Page) ([]Backup, error) {
	casted := page.(BackupPage).Body

	var resp struct {
		Backups []Backup `mapstructure:"backups" json:"backups"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Backups, err
}
