package quotasets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// QuotaSet is a set of operational limits that allow for control of compute
// usage.
type QuotaSet struct {
	// ID is tenant associated with this QuotaSet.
	ID string `json:"id"`

	// Volumes is the number of volumes that are allowed for each project.
	Volumes int `json:"volumes"`

	// Snapshots is the number of snapshots that are allowed for each project.
	Snapshots int `json:"snapshots"`

	// Gigabytes is the size (GB) of volumes and snapshots that are allowed for
	// each project.
	Gigabytes int `json:"gigabytes"`

	// PerVolumeGigabytes is the size (GB) of volumes and snapshots that are
	// allowed for each project and the specifed volume type.
	PerVolumeGigabytes int `json:"per_volume_gigabytes"`

	// Backups is the number of backups that are allowed for each project.
	Backups int `json:"backups"`

	// BackupGigabytes is the size (GB) of backups that are allowed for each
	// project.
	BackupGigabytes int `json:"backup_gigabytes"`

	// Groups is the number of groups that are allowed for each project.
	Groups int `json:"groups"`
}

// QuotaDetailSet represents details of both operational limits of compute
// resources and the current usage of those resources.
type QuotaDetailSet struct {
	// ID is the tenant ID associated with this QuotaDetailSet.
	ID string `json:"id"`

	// Volumes is the
	Volumes QuotaDetail `json:"volumes"`

	// Snapshots is the
	Snapshots QuotaDetail `json:"snapshots"`

	// Gigabytes is the
	Gigabytes QuotaDetail `json:"gigabytes"`

	// PerVolumeGigabytes is the
	PerVolumeGigabytes QuotaDetail `json:"per_volume_gigabytes"`

	// Backups is the
	Backups QuotaDetail `json:"backups"`

	// BackupGigabytes is the
	BackupGigabytes QuotaDetail `json:"backup_gigabytes"`

	// Volumes is the
	Groups QuotaDetail `json:"groups"`
}

// QuotaDetail is a set of details about a single operational limit that allows
// for control of compute usage.
type QuotaDetail struct {
	// InUse is the current number of provisioned/allocated resources of the
	// given type.
	InUse int `json:"in_use"`

	// Reserved is a transitional state when a claim against quota has been made
	// but the resource is not yet fully online.
	Reserved int `json:"reserved"`

	// Limit is the maximum number of a given resource that can be
	// allocated/provisioned.  This is what "quota" usually refers to.
	Limit int `json:"limit"`
}

// QuotaSetPage stores a single page of all QuotaSet results from a List call.
type QuotaSetPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a QuotaSetsetPage is empty.
func (page QuotaSetPage) IsEmpty() (bool, error) {
	ks, err := ExtractQuotaSets(page)
	return len(ks) == 0, err
}

// ExtractQuotaSets interprets a page of results as a slice of QuotaSets.
func ExtractQuotaSets(r pagination.Page) ([]QuotaSet, error) {
	var s struct {
		QuotaSets []QuotaSet `json:"quotas"`
	}
	err := (r.(QuotaSetPage)).ExtractInto(&s)
	return s.QuotaSets, err
}

type quotaResult struct {
	gophercloud.Result
}

// Extract is a method that attempts to interpret any QuotaSet resource response
// as a QuotaSet struct.
func (r quotaResult) Extract() (*QuotaSet, error) {
	var s struct {
		QuotaSet *QuotaSet `json:"quota_set"`
	}
	err := r.ExtractInto(&s)
	return s.QuotaSet, err
}

// GetResult is the response from a Get operation. Call its Extract method to
// interpret it as a QuotaSet.
type GetResult struct {
	quotaResult
}

// UpdateResult is the response from a Update operation. Call its Extract method
// to interpret it as a QuotaSet.
type UpdateResult struct {
	quotaResult
}

// DeleteResult is the response from a Delete operation. Call its Extract method
// to interpret it as a QuotaSet.
type DeleteResult struct {
	quotaResult
}

type quotaDetailResult struct {
	gophercloud.Result
}

// GetDetailResult is the response from a Get operation. Call its Extract
// method to interpret it as a QuotaSet.
type GetDetailResult struct {
	quotaDetailResult
}

// Extract is a method that attempts to interpret any QuotaDetailSet
// resource response as a set of QuotaDetailSet structs.
func (r quotaDetailResult) Extract() (QuotaDetailSet, error) {
	var s struct {
		QuotaData QuotaDetailSet `json:"quota_set"`
	}
	err := r.ExtractInto(&s)
	return s.QuotaData, err
}
