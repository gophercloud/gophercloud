package backups

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Backup contains all the information associated with a Cinder Backup.
type Backup struct {
	// ID is the Unique identifier of the backup.
	ID string `json:"id"`

	// CreatedAt is the date the backup was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the date the backup was updated.
	UpdatedAt time.Time `json:"-"`

	// Name is the display name of the backup.
	Name string `json:"name"`

	// Description is the description of the backup.
	Description string `json:"description"`

	// VolumeID is the ID of the Volume from which this backup was created.
	VolumeID string `json:"volume_id"`

	// SnapshotID is the ID of the snapshot from which this backup was created.
	SnapshotID string `json:"snapshot_id"`

	// Status is the status of the backup.
	Status string `json:"status"`

	// Size is the size of the backup, in GB.
	Size int `json:"size"`

	// Object Count is the number of objects in the backup.
	ObjectCount int `json:"object_count"`

	// Container is the container where the backup is stored.
	Container string `json:"container"`

	// HasDependentBackups is whether there are other backups
	// depending on this backup.
	HasDependentBackups bool `json:"has_dependent_backups"`

	// FailReason has the reason for the backup failure.
	FailReason string `json:"fail_reason"`

	// IsIncremental is whether this is an incremental backup.
	IsIncremental bool `json:"is_incremental"`

	// DataTimestamp is the time when the data on the volume was first saved.
	DataTimestamp time.Time `json:"-"`

	// ProjectID is the ID of the project that owns the backup. This is
	// an admin-only field.
	ProjectID string `json:"os-backup-project-attr:project_id"`

	// Metadata is metadata about the backup.
	// This requires microversion 3.43 or later.
	Metadata *map[string]string `json:"metadata"`

	// AvailabilityZone is the Availability Zone of the backup.
	// This requires microversion 3.51 or later.
	AvailabilityZone *string `json:"availability_zone"`
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}

// BackupPage is a pagination.Pager that is returned from a call to the List function.
type BackupPage struct {
	pagination.LinkedPageBase
}

// UnmarshalJSON converts our JSON API response into our backup struct
func (r *Backup) UnmarshalJSON(b []byte) error {
	type tmp Backup
	var s struct {
		tmp
		CreatedAt     gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt     gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
		DataTimestamp gophercloud.JSONRFC3339MilliNoZ `json:"data_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Backup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.DataTimestamp = time.Time(s.DataTimestamp)

	return err
}

// IsEmpty returns true if a BackupPage contains no Backups.
func (r BackupPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	volumes, err := ExtractBackups(r)
	return len(volumes) == 0, err
}

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

// ExtractBackups extracts and returns Backups. It is used while iterating over a backups.List call.
func ExtractBackups(r pagination.Page) ([]Backup, error) {
	var s []Backup
	err := ExtractBackupsInto(r, &s)
	return s, err
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Backup object out of the commonResult object.
func (r commonResult) Extract() (*Backup, error) {
	var s Backup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "backup")
}

func ExtractBackupsInto(r pagination.Page, v any) error {
	return r.(BackupPage).Result.ExtractIntoSlicePtr(v, "backups")
}

// RestoreResult contains the response body and error from a restore request.
type RestoreResult struct {
	commonResult
}

// Restore contains all the information associated with a Cinder Backup restore
// response.
type Restore struct {
	// BackupID is the Unique identifier of the backup.
	BackupID string `json:"backup_id"`

	// VolumeID is the Unique identifier of the volume.
	VolumeID string `json:"volume_id"`

	// Name is the name of the volume, where the backup was restored to.
	VolumeName string `json:"volume_name"`
}

// Extract will get the Backup restore object out of the RestoreResult object.
func (r RestoreResult) Extract() (*Restore, error) {
	var s Restore
	err := r.ExtractInto(&s)
	return &s, err
}

func (r RestoreResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "restore")
}

// ExportResult contains the response body and error from an export request.
type ExportResult struct {
	commonResult
}

// BackupRecord contains an information about a backup backend storage.
type BackupRecord struct {
	// The service used to perform the backup.
	BackupService string `json:"backup_service"`

	// An identifier string to locate the backup.
	BackupURL []byte `json:"backup_url"`
}

// Extract will get the Backup record object out of the ExportResult object.
func (r ExportResult) Extract() (*BackupRecord, error) {
	var s BackupRecord
	err := r.ExtractInto(&s)
	return &s, err
}

func (r ExportResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "backup-record")
}

// ImportResponse struct contains the response of the Backup Import action.
type ImportResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ImportResult contains the response body and error from an import request.
type ImportResult struct {
	gophercloud.Result
}

// Extract will get the Backup object out of the commonResult object.
func (r ImportResult) Extract() (*ImportResponse, error) {
	var s ImportResponse
	err := r.ExtractInto(&s)
	return &s, err
}

func (r ImportResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "backup")
}

// ImportBackup contains all the information to import a Cinder Backup.
type ImportBackup struct {
	ID                  string            `json:"id"`
	CreatedAt           time.Time         `json:"-"`
	UpdatedAt           time.Time         `json:"-"`
	VolumeID            string            `json:"volume_id"`
	SnapshotID          *string           `json:"snapshot_id"`
	Status              *string           `json:"status"`
	Size                *int              `json:"size"`
	ObjectCount         *int              `json:"object_count"`
	Container           *string           `json:"container"`
	ServiceMetadata     *string           `json:"service_metadata"`
	Service             *string           `json:"service"`
	Host                *string           `json:"host"`
	UserID              string            `json:"user_id"`
	DeletedAt           time.Time         `json:"-"`
	DataTimestamp       time.Time         `json:"-"`
	TempSnapshotID      *string           `json:"temp_snapshot_id"`
	TempVolumeID        *string           `json:"temp_volume_id"`
	RestoreVolumeID     *string           `json:"restore_volume_id"`
	NumDependentBackups *int              `json:"num_dependent_backups"`
	EncryptionKeyID     *string           `json:"encryption_key_id"`
	ParentID            *string           `json:"parent_id"`
	Deleted             bool              `json:"deleted"`
	DisplayName         *string           `json:"display_name"`
	DisplayDescription  *string           `json:"display_description"`
	DriverInfo          any               `json:"driver_info"`
	FailReason          *string           `json:"fail_reason"`
	ProjectID           string            `json:"project_id"`
	Metadata            map[string]string `json:"metadata"`
	AvailabilityZone    *string           `json:"availability_zone"`
}

// UnmarshalJSON converts our JSON API response into our backup struct
func (r *ImportBackup) UnmarshalJSON(b []byte) error {
	type tmp ImportBackup
	var s struct {
		tmp
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		DeletedAt     time.Time `json:"deleted_at"`
		DataTimestamp time.Time `json:"data_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ImportBackup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.DeletedAt = time.Time(s.DeletedAt)
	r.DataTimestamp = time.Time(s.DataTimestamp)

	return err
}

// MarshalJSON converts our struct request into JSON backup import request
func (r ImportBackup) MarshalJSON() ([]byte, error) {
	type b ImportBackup
	type ext struct {
		CreatedAt     *string `json:"created_at"`
		UpdatedAt     *string `json:"updated_at"`
		DeletedAt     *string `json:"deleted_at"`
		DataTimestamp *string `json:"data_timestamp"`
	}
	type tmp struct {
		b
		ext
	}

	var t ext
	if r.CreatedAt != (time.Time{}) {
		v := r.CreatedAt.Format(time.RFC3339)
		t.CreatedAt = &v
	}
	if r.UpdatedAt != (time.Time{}) {
		v := r.UpdatedAt.Format(time.RFC3339)
		t.UpdatedAt = &v
	}
	if r.DeletedAt != (time.Time{}) {
		v := r.DeletedAt.Format(time.RFC3339)
		t.DeletedAt = &v
	}
	if r.DataTimestamp != (time.Time{}) {
		v := r.DataTimestamp.Format(time.RFC3339)
		t.DataTimestamp = &v
	}

	if r.Metadata == nil {
		r.Metadata = make(map[string]string)
	}

	s := tmp{
		b(r),
		t,
	}

	return json.Marshal(s)
}

// ResetStatusResult contains the response error from a ResetStatus request.
type ResetStatusResult struct {
	gophercloud.ErrResult
}

// ForceDeleteResult contains the response error from a ForceDelete request.
type ForceDeleteResult struct {
	gophercloud.ErrResult
}
