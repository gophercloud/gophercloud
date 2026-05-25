package backups

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Backup represents a database backup.
type Backup struct {
	Created       time.Time `json:"-"`
	Datastore     datastores.DatastorePartial
	Description   string
	ID            string
	InstanceID    string `json:"instance_id"`
	LocationRef   string `json:"locationRef"`
	Name          string
	ParentID      string `json:"parent_id"`
	ProjectID     string `json:"project_id"`
	Size          float64
	Status        string
	StorageDriver string    `json:"storage_driver"`
	Updated       time.Time `json:"-"`
}

// UnmarshalJSON converts backup timestamps to time.Time.
func (r *Backup) UnmarshalJSON(b []byte) error {
	type tmp Backup
	var s struct {
		tmp
		Created gophercloud.JSONRFC3339NoZ `json:"created"`
		Updated gophercloud.JSONRFC3339NoZ `json:"updated"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Backup(s.tmp)
	r.Created = time.Time(s.Created)
	r.Updated = time.Time(s.Updated)
	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract retrieves a Backup resource from an operation result.
func (r commonResult) Extract() (*Backup, error) {
	var s struct {
		Backup *Backup `json:"backup"`
	}
	err := r.ExtractInto(&s)
	return s.Backup, err
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// BackupPage represents a page of database backups.
type BackupPage struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a BackupPage is empty.
func (r BackupPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	backups, err := ExtractBackups(r)
	return len(backups) == 0, err
}

// ExtractBackups retrieves a slice of Backup structs from a paginated
// collection.
func ExtractBackups(r pagination.Page) ([]Backup, error) {
	var s struct {
		Backups []Backup `json:"backups"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Backups, err
}
