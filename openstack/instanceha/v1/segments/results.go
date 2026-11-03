package segments

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Segment represents a Masakari segment.
type Segment struct {
	// ID is the numeric database identifier used for pagination markers.
	ID int `json:"id"`
	// UUID is the unique identifier used in detail requests.
	UUID string `json:"uuid"`
	// Name is the failover segment name.
	Name string `json:"name"`
	// Description describes the failover segment.
	Description string `json:"description"`
	// ServiceType is the service managed by the segment, such as COMPUTE.
	ServiceType string `json:"service_type"`
	// RecoveryMethod is the recovery strategy: auto, reserved_host, auto_priority, or rh_priority.
	RecoveryMethod string `json:"recovery_method"`
	// Enabled indicates whether recovery is enabled for the segment. Requires microversion 1.2 or later.
	Enabled *bool `json:"enabled"`
	// CreatedAt is the creation time. A null timestamp is decoded as the zero time.
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	// UpdatedAt is the last update time. A null timestamp is decoded as the zero time.
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	// DeletedAt is the deletion time. A null timestamp is decoded as the zero time.
	DeletedAt gophercloud.JSONRFC3339MilliNoZ `json:"deleted_at"`
	// Deleted indicates whether the resource has been soft-deleted.
	Deleted bool `json:"deleted"`
}

// SegmentPage is a page of failover segments returned by List.
type SegmentPage struct{ pagination.MarkerPageBase }

// ResourceKey identifies the collection in the JSON response.
func (page SegmentPage) ResourceKey() string { return "segments" }

// IsEmpty reports whether the page contains any resources.
func (page SegmentPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractSegments(page)
	return len(items) == 0, err
}

// LastMarker returns the last numeric ID for the next pagination request.
// Missing or non-advancing markers are returned as errors.
func (page SegmentPage) LastMarker() (string, error) {
	items, err := ExtractSegments(page)
	if err != nil || len(items) == 0 {
		return "", err
	}
	id := items[len(items)-1].ID
	if id <= 0 {
		return "", fmt.Errorf("missing numeric marker in segments response")
	}
	if strconv.Itoa(id) == page.URL.Query().Get("marker") {
		return "", fmt.Errorf("non-advancing segments marker: %d", id)
	}
	return strconv.Itoa(id), nil
}

// ExtractSegments extracts failover segments from a page.
func ExtractSegments(r pagination.Page) ([]Segment, error) {
	var s struct {
		Items []Segment `json:"segments"`
	}
	var err error
	switch page := r.(type) {
	case *SegmentPage:
		err = page.ExtractInto(&s)
	case SegmentPage:
		err = page.ExtractInto(&s)
	default:
		return nil, fmt.Errorf("unexpected segments page type %T", r)
	}
	return s.Items, err
}

// GetResult contains the response to a Get request. Call Extract to retrieve the segment.
type GetResult struct{ gophercloud.Result }

// Extract decodes the segment returned by Get.
func (r GetResult) Extract() (*Segment, error) {
	var s struct {
		Item *Segment `json:"segment"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
