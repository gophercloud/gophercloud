package backup

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/openstack/db/v1/datastores"
)

// Database represents a Backup API resource.
type Backup struct {
	// The ID of the backup.
	ID string `json:"id"`
	// Name of the backup.
	Name string `json:"name"`
	// The ID of the instance to create backup for.
	InstanceID string `json:"instance_id"`
	// A datastore object.
	Datastore datastores.DatastorePartial `json:"datastore"`
	// The URL of the backup location.
	LocationRef string `json:"locationRef"`
	// ID of the parent backup to perform an incremental backup from.
	ParentID string `json:"parent_id"`
	// Size of the backup, the unit is GB.
	Size float64 `json:"size"`
	// Status of the backup.
	Status string `json:"status"`
	// The date and time when the resource was created.The date and time
	// stamp format is ISO 8601: CCYY-MM-DDThh:mm:ss±hh:mm
	//	For example, 2015-08-27T09:49:58-05:00.
	//	The ±hh:mm value, if included, is the time zone as an offset from UTC.
	// In the previous example, the offset value is -05:00.
	Created string `json:"created"`
	// The date and time when the resource was updated. The date
	// and time stamp format is ISO 8601:"CCYY-MM-DDThh:mm:ss±hh:mm"
	// The ±hh:mm value, if included, is the time zone as an offset from UTC.
	// For example, 2015-08-27T09:49:58-05:00.The UTC time zone is assumed.
	Updated string `json:"updated"`
	// An optional description for the backup.
	Description string `json:"description"`
}

// CreateResult represents the result of a Create operation.

type backupBaseResult struct {
	gophercloud.Result
}
type CreateResult struct {
	backupBaseResult
}

type GetResult struct {
	backupBaseResult
}

func (r backupBaseResult) Extract() (*Backup, error) {
	var b struct {
		Backup *Backup `json:"backup"`
	}
	err := r.ExtractInto(&b)
	return b.Backup, err
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// BackupPage represents a single page of a paginated Backup collection.
type BackupsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page BackupsPage) IsEmpty() (bool, error) {
	backups, err := ExtractBackups(page)
	return len(backups) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page BackupsPage) NextPageURL() (string, error) {
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
func ExtractBackups(page pagination.Page) ([]Backup, error) {
	r := page.(BackupsPage)
	var s struct {
		Backups []Backup `json:"backups"`
	}
	err := r.ExtractInto(&s)
	return s.Backups, err
}
