package hosts

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/segments"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Host represents a Masakari host.
type Host struct {
	// ID is the numeric database identifier used for pagination markers.
	ID int `json:"id"`
	// UUID is the unique identifier used in detail requests.
	UUID string `json:"uuid"`
	// Name is the host name.
	Name string `json:"name"`
	// Type is the host type, such as COMPUTE_HOST.
	Type string `json:"type"`
	// ControlAttributes describes the host access mechanism, such as SSH.
	ControlAttributes string `json:"control_attributes"`
	// Reserved indicates whether the host is reserved for recovery.
	Reserved bool `json:"reserved"`
	// OnMaintenance indicates whether the host is in maintenance mode.
	OnMaintenance bool `json:"on_maintenance"`
	// FailoverSegmentID is the UUID of the segment containing the host.
	FailoverSegmentID string `json:"failover_segment_id"`
	// FailoverSegment contains the segment details embedded in the host response.
	FailoverSegment *segments.Segment `json:"failover_segment"`
	// CreatedAt is the creation time. A null timestamp is decoded as the zero time.
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	// UpdatedAt is the last update time. A null timestamp is decoded as the zero time.
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	// DeletedAt is the deletion time. A null timestamp is decoded as the zero time.
	DeletedAt gophercloud.JSONRFC3339MilliNoZ `json:"deleted_at"`
	// Deleted indicates whether the resource has been soft-deleted.
	Deleted bool `json:"deleted"`
}

// HostPage is a page of hosts belonging to a failover segment returned by List.
type HostPage struct{ pagination.MarkerPageBase }

// ResourceKey identifies the collection in the JSON response.
func (page HostPage) ResourceKey() string { return "hosts" }

// IsEmpty reports whether the page contains any resources.
func (page HostPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractHosts(page)
	return len(items) == 0, err
}

// LastMarker returns the last numeric ID for the next pagination request.
// Missing or non-advancing markers are returned as errors.
func (page HostPage) LastMarker() (string, error) {
	items, err := ExtractHosts(page)
	if err != nil || len(items) == 0 {
		return "", err
	}
	id := items[len(items)-1].ID
	if id <= 0 {
		return "", fmt.Errorf("missing numeric marker in hosts response")
	}
	if strconv.Itoa(id) == page.URL.Query().Get("marker") {
		return "", fmt.Errorf("non-advancing hosts marker: %d", id)
	}
	return strconv.Itoa(id), nil
}

// ExtractHosts extracts hosts belonging to a failover segment from a page.
func ExtractHosts(r pagination.Page) ([]Host, error) {
	var s struct {
		Items []Host `json:"hosts"`
	}
	var err error
	switch page := r.(type) {
	case *HostPage:
		err = page.ExtractInto(&s)
	case HostPage:
		err = page.ExtractInto(&s)
	default:
		return nil, fmt.Errorf("unexpected hosts page type %T", r)
	}
	return s.Items, err
}

// GetResult contains the response to a Get request. Call Extract to retrieve the host.
type GetResult struct{ gophercloud.Result }

// Extract decodes the host returned by Get.
func (r GetResult) Extract() (*Host, error) {
	var s struct {
		Item *Host `json:"host"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
