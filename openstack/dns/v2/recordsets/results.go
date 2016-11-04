package recordsets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult, CreateResult or UpdateResult as a concrete Zone.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (RecordSet, error) {
	var s RecordSet
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult is the deferred result of a Create call.
type CreateResult struct {
	commonResult
}

// GetResult is the deferred result of a Get call.
type GetResult struct {
	commonResult
}

// UpdateResult is the deferred result of an Update call.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the deferred result of an Delete call.
type DeleteResult struct {
	gophercloud.ErrResult
}

// RRSetPage is a single page of RecordSet results.
type RRSetPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the page contains no results.
func (p RRSetPage) IsEmpty() (bool, error) {
	services, err := ExtractRRSets(p)
	return len(services) == 0, err
}

// ExtractServices extracts a slice of Services from a Collection acquired from List.
func ExtractRRSets(r pagination.Page) ([]RecordSet, error) {
	var s struct {
		RRSets []RecordSet `json:"recordsets"`
	}
	err := (r.(RRSetPage)).ExtractInto(&s)
	return s.RRSets, err
}

type RecordSet struct {
	ID        string                          `json:"id,omitempty"`
	ZoneID    string                          `json:"zone_id,omitempty"`
	ProjectID string                          `json:"project_id,omitempty"`
	Name      string                          `json:"name,omitempty"`
	ZoneName  string                          `json:"zone_name,omitempty"`
	Type      string                          `json:"type,omitempty"`
	Records   []string                        `json:"records,omitempty"`
	TTL       int                             `json:"ttl,omitempty"`
	Status    string                          `json:"status,omitempty"`
	Action    string                          `json:"action,omitempty"`
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at,omitempty"`
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at,omitempty"`
}
