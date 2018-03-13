package apiversion

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"time"
)

type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations. Call its Extract method to
// interpret it as a Version.
type GetResult struct {
	commonResult
}

// Version represents an API Version
type Version struct {
	// ID is the Version's unique ID.
	ID         string                   `json:"id"`
	Links      []gophercloud.Link       `json:"links"`
	MaxVersion string                   `json:"max_version"`
	MediaTypes []map[string]interface{} `json:"media-types"`
	MinVersion string                   `json:"min_version"`
	Status     string                   `json:"status"`
	Updated    time.Time                `json:"updated"`
}

// VersionPage contains a single page of all Versions from a ListDetails call.
type VersionPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a VersionPage contains any results.
func (page VersionPage) IsEmpty() (bool, error) {
	versions, err := ExtractVersions(page)
	return len(versions) == 0, err
}

// Extract provides access to the individual event returned by Get and extracts Event
func (r commonResult) Extract() (*Version, error) {
	var s struct {
		Version *Version `json:"version"`
	}
	err := r.ExtractInto(&s)
	return s.Version, err
}

// ExtractVersions provides access to the list of version in a page acquired from the ListDetail operation.
func ExtractVersions(r pagination.Page) ([]Version, error) {
	var s struct {
		Versions []Version `json:"versions"`
	}
	err := (r.(VersionPage)).ExtractInto(&s)
	return s.Versions, err
}
