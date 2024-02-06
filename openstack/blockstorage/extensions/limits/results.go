package limits

import (
	"github.com/gophercloud/gophercloud/v2"
)

// Limits is a struct that contains the response of a limit query.
type Limits struct {
	// Absolute contains the limits and usage information.
	// An absolute limit value of -1 indicates that the absolute limit for the item is infinite.
	Absolute Absolute `json:"absolute"`
	// Rate contains rate-limit volume copy bandwidth, used to mitigate slow down of data access from the instances.
	Rate []Rate `json:"rate"`
}

// Absolute is a struct that contains the current resource usage and limits
// of a project.
type Absolute struct {
	// MaxTotalVolumes is the maximum number of volumes.
	MaxTotalVolumes int `json:"maxTotalVolumes"`

	// MaxTotalSnapshots is the maximum number of snapshots.
	MaxTotalSnapshots int `json:"maxTotalSnapshots"`

	// MaxTotalVolumeGigabytes is the maximum total amount of volumes, in gibibytes (GiB).
	MaxTotalVolumeGigabytes int `json:"maxTotalVolumeGigabytes"`

	// MaxTotalBackups is the maximum number of backups.
	MaxTotalBackups int `json:"maxTotalBackups"`

	// MaxTotalBackupGigabytes is the maximum total amount of backups, in gibibytes (GiB).
	MaxTotalBackupGigabytes int `json:"maxTotalBackupGigabytes"`

	// TotalVolumesUsed is the total number of volumes used.
	TotalVolumesUsed int `json:"totalVolumesUsed"`

	// TotalGigabytesUsed is the total number of gibibytes (GiB) used.
	TotalGigabytesUsed int `json:"totalGigabytesUsed"`

	// TotalSnapshotsUsed the total number of snapshots used.
	TotalSnapshotsUsed int `json:"totalSnapshotsUsed"`

	// TotalBackupsUsed is the total number of backups used.
	TotalBackupsUsed int `json:"totalBackupsUsed"`

	// TotalBackupGigabytesUsed is the total number of backups gibibytes (GiB) used.
	TotalBackupGigabytesUsed int `json:"totalBackupGigabytesUsed"`
}

// Rate is a struct that contains the
// rate-limit volume copy bandwidth, used to mitigate slow down of data access from the instances.
type Rate struct {
	Regex string  `json:"regex"`
	URI   string  `json:"uri"`
	Limit []Limit `json:"limit"`
}

// Limit struct contains Limit values for the Rate struct
type Limit struct {
	Verb          string `json:"verb"`
	NextAvailable string `json:"next-available"`
	Unit          string `json:"unit"`
	Value         int    `json:"value"`
	Remaining     int    `json:"remaining"`
}

// Extract interprets a limits result as a Limits.
func (r GetResult) Extract() (*Limits, error) {
	var s struct {
		Limits *Limits `json:"limits"`
	}
	err := r.ExtractInto(&s)
	return s.Limits, err
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as an Absolute.
type GetResult struct {
	gophercloud.Result
}
