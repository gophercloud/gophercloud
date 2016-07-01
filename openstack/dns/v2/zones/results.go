package zones

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult, CreateResult or UpdateResult as a concrete Zone.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (Zone, error) {
	var s Zone
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

// ZonePage is a single page of Zonje results.
type ZonePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the page contains no results.
func (p ZonePage) IsEmpty() (bool, error) {
	services, err := ExtractZones(p)
	return len(services) == 0, err
}

// ExtractServices extracts a slice of Services from a Collection acquired from List.
func ExtractZones(r pagination.Page) ([]Zone, error) {
	var s struct {
		Zones []Zone `json:"zones"`
	}
	err := (r.(ZonePage)).ExtractInto(&s)
	return s.Zones, err
}

type Zone struct {
	ProjectID string  `json:"project_id,omitempty"`
	ID        string  `json:"id,omitempty"`
	PoolID    string  `json:"pool_id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Email     string  `json:"email,omitempty"`
	TTL       int     `json:"ttl,omitempty"`
	Status    string  `json:"status,omitempty"`
	Serial    float64 `json:"serial,omitempty"`
	Type      string  `json:"type,omitempty"`
}
