package notifications

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Notification struct {
	ID                  int                             `json:"id"`
	UUID                string                          `json:"notification_uuid"`
	SourceHostUUID      string                          `json:"source_host_uuid"`
	FailoverSegmentUUID string                          `json:"failover_segment_uuid"`
	Type                string                          `json:"type"`
	Status              string                          `json:"status"`
	GeneratedTime       gophercloud.JSONRFC3339MilliNoZ `json:"generated_time"`
	CreatedAt           gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	UpdatedAt           gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	Payload             map[string]any                  `json:"payload"`
	Message             string                          `json:"message"`
}

type NotificationPage struct{ pagination.MarkerPageBase }

func (page NotificationPage) ResourceKey() string { return "notifications" }

func (page NotificationPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractNotifications(page)
	return len(items) == 0, err
}

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

type GetResult struct{ gophercloud.Result }

func (r GetResult) Extract() (*Notification, error) {
	var s struct {
		Item *Notification `json:"notification"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
