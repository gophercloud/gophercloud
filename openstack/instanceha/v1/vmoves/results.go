package vmoves

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type VMove struct {
	ID               int                             `json:"id"`
	UUID             string                          `json:"uuid"`
	NotificationUUID string                          `json:"notification_uuid"`
	InstanceUUID     string                          `json:"instance_uuid"`
	InstanceName     string                          `json:"instance_name"`
	SourceHost       string                          `json:"source_host"`
	DestHost         string                          `json:"dest_host"`
	Type             string                          `json:"type"`
	Status           string                          `json:"status"`
	StartTime        gophercloud.JSONRFC3339MilliNoZ `json:"start_time"`
	EndTime          gophercloud.JSONRFC3339MilliNoZ `json:"end_time"`
	Message          string                          `json:"message"`
}

type VMovePage struct{ pagination.MarkerPageBase }

func (page VMovePage) ResourceKey() string { return "vmoves" }

func (page VMovePage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractVMoves(page)
	return len(items) == 0, err
}

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

type GetResult struct{ gophercloud.Result }

func (r GetResult) Extract() (*VMove, error) {
	var s struct {
		Item *VMove `json:"vmove"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
