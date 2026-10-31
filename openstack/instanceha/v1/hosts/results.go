package hosts

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Host struct {
	ID                int    `json:"id"`
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	ControlAttributes string `json:"control_attributes"`
	Reserved          bool   `json:"reserved"`
	OnMaintenance     bool   `json:"on_maintenance"`
	FailoverSegmentID string `json:"failover_segment_id"`
}

type HostPage struct{ pagination.MarkerPageBase }

func (page HostPage) ResourceKey() string { return "hosts" }

func (page HostPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractHosts(page)
	return len(items) == 0, err
}

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

type GetResult struct{ gophercloud.Result }

func (r GetResult) Extract() (*Host, error) {
	var s struct {
		Item *Host `json:"host"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
