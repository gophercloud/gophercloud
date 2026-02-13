package quotasets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// QuotaSet is a set of operational limits that allow for control of manila
// usage.
type QuotaSet struct {
	// Gigabytes is the total size of share storage for the project in gigabytes.
	Gigabytes *int `json:"gigabytes,omitempty"`

	// Snapshots is the total number of share snapshots for the project.
	Snapshots *int `json:"snapshots,omitempty"`

	// Shares is the total number of shares for the project.
	Shares *int `json:"shares,omitempty"`

	// SnapshotGigabytes is the total size of share snapshots for the project in gigabytes.
	SnapshotGigabytes *int `json:"snapshot_gigabytes,omitempty"`

	// Share network is the total number of share networks for the project.
	ShareNetworks *int `json:"share_networks,omitempty"`

	// Share groups is the total number of share groups for the project.
	ShareGroups *int `json:"share_groups,omitempty"`

	// Share group snapshots is the total number of share group snapshots for the project.
	ShareGroupSnapshots *int `json:"share_group_snapshots,omitempty"`

	// Share Replicas is the total number of share replicas for the project.
	ShareReplicas *int `json:"share_replicas,omitempty"`

	// Share Replica Gigabytes is the total size of share replicas for the project in gigabytes.
	ShareReplicaGigabytes *int `json:"share_replica_gigabytes,omitempty"`

	// PerShareGigabytes is the maximum size of a share for the project in gigabytes.
	PerShareGigabytes *int `json:"per_share_gigabytes,omitempty"`
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
