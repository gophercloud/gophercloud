package quotasets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// QuotaSet is a set of operational limits that allow for control of compute
// usage.
type QuotaSet struct {
	// ID is project associated with this QuotaSet.
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
	// ID is the project ID associated with this QuotaDetailSet.
	ID string `json:"id"`

	// Volumes is the volume usage information for this project, including
	// in_use, limit, reserved and allocated attributes. Note: allocated
	// attribute is available only when nested quota is enabled.
	Volumes QuotaDetail `json:"volumes"`

	// Snapshots is the snapshot usage information for this project, including
	// in_use, limit, reserved and allocated attributes. Note: allocated
	// attribute is available only when nested quota is enabled.
	Snapshots QuotaDetail `json:"snapshots"`

	// Gigabytes is the size (GB) usage information of volumes and snapshots
	// for this project, including in_use, limit, reserved and allocated
	// attributes. Note: allocated attribute is available only when nested
	// quota is enabled.
	Gigabytes QuotaDetail `json:"gigabytes"`

	// PerVolumeGigabytes is the size (GB) usage information for each volume,
	// including in_use, limit, reserved and allocated attributes. Note:
	// allocated attribute is available only when nested quota is enabled and
	// only limit is meaningful here.
	PerVolumeGigabytes QuotaDetail `json:"per_volume_gigabytes"`

	// Backups is the backup usage information for this project, including
	// in_use, limit, reserved and allocated attributes. Note: allocated
	// attribute is available only when nested quota is enabled.
	Backups QuotaDetail `json:"backups"`

	// BackupGigabytes is the size (GB) usage information of backup for this
	// project, including in_use, limit, reserved and allocated attributes.
	// Note: allocated attribute is available only when nested quota is
	// enabled.
	BackupGigabytes QuotaDetail `json:"backup_gigabytes"`

	// Volumes is the group usage information for this project, including
	// in_use, limit, reserved and allocated attributes. Note: allocated
	// attribute is available only when nested quota is enabled.
	Groups QuotaDetail `json:"groups"`
}

// QuotaDetail is a set of details about a single operational limit that allows
// for control of compute usage.
type QuotaDetail struct {
	// InUse is the current number of provisioned resources of the given type.
	InUse int `json:"in_use"`

	// Allocated is the current number of resources of a given type allocated
	// for use.  It is only available when nested quota is enabled. It tells
	// how
	Allocated int `json:"allocated"`

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
