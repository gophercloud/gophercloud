package snapshots

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Snapshot contains all the information associated with a Cinder Snapshot.
type Snapshot struct {
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

	// ID of the Volume from which this Snapshot was created.
	VolumeID string `json:"volume_id"`

	// Currect status of the Snapshot.
	Status string `json:"status"`

	// Size of the Snapshot, in GB.
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

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// SnapshotPage is a pagination.Pager that is returned from a call to the List function.
type SnapshotPage struct {
	pagination.LinkedPageBase
}

// UnmarshalJSON converts our JSON API response into our snapshot struct
func (r *Snapshot) UnmarshalJSON(b []byte) error {
	type tmp Snapshot
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Snapshot(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

// IsEmpty returns true if a SnapshotPage contains no Snapshots.
func (r SnapshotPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	volumes, err := ExtractSnapshots(r)
	return len(volumes) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r SnapshotPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"snapshots_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractSnapshotsInto similar to ExtractInto but operates on a `list` of snapshots
func ExtractSnapshotsInto(r pagination.Page, v interface{}) error {
	return r.(SnapshotPage).Result.ExtractIntoSlicePtr(v, "snapshots")
}

// ExtractSnapshots extracts and returns Snapshots. It is used while iterating over a snapshots.List call.
func ExtractSnapshots(r pagination.Page) ([]Snapshot, error) {
	var s []Snapshot
	err := ExtractSnapshotsInto(r, &s)
	return s, err
}

// UpdateMetadataResult contains the response body and error from an UpdateMetadata request.
type UpdateMetadataResult struct {
	commonResult
}

// ExtractMetadata returns the metadata from a response from snapshots.UpdateMetadata.
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

// Extract will get the Snapshot object out of the commonResult object.
func (r commonResult) Extract() (*Snapshot, error) {
	var s Snapshot
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a snapshot struct
func (r commonResult) ExtractInto(s interface{}) error {
	return r.Result.ExtractIntoStructPtr(s, "snapshot")
}

// ResetStatusResult contains the response error from a ResetStatus request.
type ResetStatusResult struct {
	gophercloud.ErrResult
}

// UpdateStatusResult contains the response error from an UpdateStatus request.
type UpdateStatusResult struct {
	gophercloud.ErrResult
}

// ForceDeleteResult contains the response error from a ForceDelete request.
type ForceDeleteResult struct {
	gophercloud.ErrResult
}
