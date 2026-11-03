package notifications

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Notification represents a Masakari notification.
type Notification struct {
	// ID is the numeric database identifier used for pagination markers.
	ID int `json:"id"`
	// UUID is the unique identifier used in detail requests.
	UUID string `json:"notification_uuid"`
	// SourceHostUUID is the UUID of the host that generated the notification.
	SourceHostUUID string `json:"source_host_uuid"`
	// FailoverSegmentUUID is the UUID of the affected segment.
	FailoverSegmentUUID string `json:"failover_segment_uuid"`
	// Type is the event type, such as COMPUTE_HOST, PROCESS, or VM.
	Type string `json:"type"`
	// Status is the resource processing status.
	Status string `json:"status"`
	// GeneratedTime is when the notification was generated.
	GeneratedTime gophercloud.JSONRFC3339MilliNoZ `json:"generated_time"`
	// CreatedAt is the creation time. A null timestamp is decoded as the zero time.
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	// UpdatedAt is the last update time. A null timestamp is decoded as the zero time.
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	// Payload contains the notification event details.
	Payload map[string]any `json:"payload"`
	// Message contains additional processing information.
	Message string `json:"message"`
}

// NotificationPage is a page of failure notifications returned by List.
type NotificationPage struct{ pagination.MarkerPageBase }

// ResourceKey identifies the collection in the JSON response.
func (page NotificationPage) ResourceKey() string { return "notifications" }

// IsEmpty reports whether the page contains any resources.
func (page NotificationPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractNotifications(page)
	return len(items) == 0, err
}

// LastMarker returns the last numeric ID for the next pagination request.
// Missing or non-advancing markers are returned as errors.
func (page NotificationPage) LastMarker() (string, error) {
	items, err := ExtractNotifications(page)
	if err != nil || len(items) == 0 {
		return "", err
	}
	id := items[len(items)-1].ID
	if id <= 0 {
		return "", fmt.Errorf("missing numeric marker in notifications response")
	}
	if strconv.Itoa(id) == page.URL.Query().Get("marker") {
		return "", fmt.Errorf("non-advancing notifications marker: %d", id)
	}
	return strconv.Itoa(id), nil
}

// ExtractNotifications extracts failure notifications from a page.
func ExtractNotifications(r pagination.Page) ([]Notification, error) {
	var s struct {
		Items []Notification `json:"notifications"`
	}
	var err error
	switch page := r.(type) {
	case *NotificationPage:
		err = page.ExtractInto(&s)
	case NotificationPage:
		err = page.ExtractInto(&s)
	default:
		return nil, fmt.Errorf("unexpected notifications page type %T", r)
	}
	return s.Items, err
}

// GetResult contains the response to a Get request. Call Extract to retrieve the notification.
type GetResult struct{ gophercloud.Result }

// Extract decodes the notification returned by Get.
func (r GetResult) Extract() (*Notification, error) {
	var s struct {
		Item *Notification `json:"notification"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
