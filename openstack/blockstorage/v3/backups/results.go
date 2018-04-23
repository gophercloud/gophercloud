package backups

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Backup contains all the information associated with a Cinder Backup.
type Backup struct {
	// Unique identifier.
	ID string `json:"id"`

	// Date created.
	CreatedAt time.Time `json:"-"`

	// Date updated.
	UpdatedAt time.Time `json:"-"`

	// Display name.
	Name string `json:"name"`

	// Display description.
	Description string `json:"description"`

	// ID of the Volume from which this Backup was created.
	VolumeID string `json:"volume_id"`

	// Currect status of the Backup.
	Status string `json:"status"`

	// Size of the Backup, in GB.
	Size int `json:"size"`

	// User-defined key-value pairs.
	Metadata map[string]string `json:"metadata"`
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
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Backup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

// IsEmpty returns true if a BackupPage contains no Backups.
func (r BackupPage) IsEmpty() (bool, error) {
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
	var s struct {
		Backups []Backup `json:"backups"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Backups, err
}

// UpdateMetadataResult contains the response body and error from an UpdateMetadata request.
type UpdateMetadataResult struct {
	commonResult
}

// ExtractMetadata returns the metadata from a response from backups.UpdateMetadata.
func (r UpdateMetadataResult) ExtractMetadata() (map[string]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	m := r.Body.(map[string]interface{})["metadata"]
	return m.(map[string]interface{}), nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Backup object out of the commonResult object.
func (r commonResult) Extract() (*Backup, error) {
	var s struct {
		Backup *Backup `json:"backup"`
	}
	err := r.ExtractInto(&s)
	return s.Backup, err
}
