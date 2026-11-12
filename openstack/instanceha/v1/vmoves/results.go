package vmoves

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// VMove represents a Masakari VM move. Requires microversion 1.3 or later.
type VMove struct {
	// ID is the numeric database identifier used for pagination markers.
	ID int `json:"id"`
	// UUID is the unique identifier used in detail requests.
	UUID string `json:"uuid"`
	// NotificationUUID is the UUID of the parent COMPUTE_HOST notification.
	NotificationUUID string `json:"notification_uuid"`
	// InstanceUUID is the UUID of the affected server.
	InstanceUUID string `json:"instance_uuid"`
	// InstanceName is the name of the affected server.
	InstanceName string `json:"instance_name"`
	// SourceHost is the source compute host.
	SourceHost string `json:"source_host"`
	// DestHost is the destination compute host.
	DestHost string `json:"dest_host"`
	// Type is the VM move type, such as evacuation or migration.
	Type string `json:"type"`
	// Status is the resource processing status.
	Status string `json:"status"`
	// StartTime is when the VM move started. A null timestamp is decoded as the zero time.
	StartTime gophercloud.JSONRFC3339MilliNoZ `json:"start_time"`
	// EndTime is when the VM move finished. A null timestamp is decoded as the zero time.
	EndTime gophercloud.JSONRFC3339MilliNoZ `json:"end_time"`
	// Message contains additional processing information.
	Message string `json:"message"`
}

// VMovePage is a page of VM moves belonging to a COMPUTE_HOST notification returned by List.
type VMovePage struct{ pagination.MarkerPageBase }

// ResourceKey identifies the collection in the JSON response.
func (page VMovePage) ResourceKey() string { return "vmoves" }

// IsEmpty reports whether the page contains any resources.
func (page VMovePage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractVMoves(page)
	return len(items) == 0, err
}

// LastMarker returns the last numeric ID for the next pagination request.
// Missing or non-advancing markers are returned as errors.
func (page VMovePage) LastMarker() (string, error) {
	items, err := ExtractVMoves(page)
	if err != nil || len(items) == 0 {
		return "", err
	}
	id := items[len(items)-1].ID
	if id <= 0 {
		return "", fmt.Errorf("missing numeric marker in vmoves response")
	}
	if strconv.Itoa(id) == page.URL.Query().Get("marker") {
		return "", fmt.Errorf("non-advancing vmoves marker: %d", id)
	}
	return strconv.Itoa(id), nil
}

// ExtractVMoves extracts VM moves belonging to a COMPUTE_HOST notification from a page.
func ExtractVMoves(r pagination.Page) ([]VMove, error) {
	var s struct {
		Items []VMove `json:"vmoves"`
	}
	var err error
	switch page := r.(type) {
	case *VMovePage:
		err = page.ExtractInto(&s)
	case VMovePage:
		err = page.ExtractInto(&s)
	default:
		return nil, fmt.Errorf("unexpected vmoves page type %T", r)
	}
	return s.Items, err
}

// GetResult contains the response to a Get request. Call Extract to retrieve the VM move.
type GetResult struct{ gophercloud.Result }

// Extract decodes the VM move returned by Get.
func (r GetResult) Extract() (*VMove, error) {
	var s struct {
		Item *VMove `json:"vmove"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
