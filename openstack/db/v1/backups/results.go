package backups

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/pagination"
)

// Backup represents a database backup
type Backup struct {
	// Indicates the datetime that the backup was created
	Created time.Time `json:"created"`

	// Indicates the most recent datetime that the backup was updated.
	Updated time.Time `json:"updated"`

	// Indicates the optional description for the backup.
	Description string `json:"description"`

	// Indicates the unique identifier for the backup resource.
	ID string `json:"id"`

	// Indicates how the instance stores data.
	Datastore datastores.DatastorePartial `json:"datastore"`

	// Indicates the unique identifier of the instance to create backup for.
	InstanceID string `json:"instance_id"`

	// Indicates the URL of the backup location
	LocationRef string `json:"locationRef"`

	// The human-readable name of the instance.
	Name string `json:"name"`

	// Indicates the unique identifier of the parent backup to perform
	// an incremental backup from.
	ParentId string `json:"parent_id"`

	// Indicates the volume size of the backup in gigabytes (GB)
	Size float64 `json:"size"`

	// The build status of the backup.
	Status string `json:"status"`

	// Indicates the unique identifier of the project ID of the backup
	ProjectId string `json:"project_id"`
}

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

// Extract will extract an Backup from various result structs.
func (r commonResult) Extract() (*Backup, error) {
	var s struct {
		Backup *Backup `json:"backup"`
	}
	err := r.ExtractInto(&s)
	return s.Backup, err
}

// BackupPage represents a single page of a paginated user collection.
type BackupPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page BackupPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	users, err := ExtractBackups(page)
	return len(users) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page BackupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"backups_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractBackups will convert a generic pagination struct into a more
// relevant slice of Backup structs.
func ExtractBackups(r pagination.Page) ([]Backup, error) {
	var s struct {
		Backups []Backup `json:"backups"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Backups, err
}
