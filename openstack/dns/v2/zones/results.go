package zones

import (
	"encoding/json"
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

// GetResult is the deferred result of a Get call.
type GetResult struct {
	commonResult
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

// ExtractZones extracts a slice of Services from a Collection acquired from List.
func ExtractZones(r pagination.Page) ([]Zone, error) {
	var s struct {
		Zones []Zone `json:"zones"`
	}
	err := (r.(ZonePage)).ExtractInto(&s)
	return s.Zones, err
}

type Zone struct {
	ID            string                          `json:"id,omitempty"`
	PoolID        string                          `json:"pool_id,omitempty"`
	ProjectID     string                          `json:"project_id,omitempty"`
	Name          string                          `json:"name,omitempty"`
	Email         string                          `json:"email,omitempty"`
	TTL           int                             `json:"ttl,omitempty"`
	Serial        json.Number                     `json:"serial,omitempty"`
	Status        string                          `json:"status,omitempty"`
	Action        string                          `json:"action,omitempty"`
	Version       int                             `json:"version,omitempty"`
	Attributes    map[string]string               `json:"attributes,omitempty"`
	Type          string                          `json:"type,omitempty"`
	Masters       []string                        `json:"masters,omitempty"`
	CreatedAt     gophercloud.JSONRFC3339MilliNoZ `json:"created_at,omitempty"`
	UpdatedAt     gophercloud.JSONRFC3339MilliNoZ `json:"updated_at,omitempty"`
	TransferredAt gophercloud.JSONRFC3339MilliNoZ `json:"transferred_at,omitempty"`
}
