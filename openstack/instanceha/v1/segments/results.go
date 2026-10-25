package segments

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Segment struct {
	ID             int    `json:"id"`
	UUID           string `json:"uuid"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ServiceType    string `json:"service_type"`
	RecoveryMethod string `json:"recovery_method"`
	Enabled        *bool  `json:"enabled"`
}

type SegmentPage struct{ pagination.MarkerPageBase }

func (page SegmentPage) ResourceKey() string { return "segments" }

func (page SegmentPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	items, err := ExtractSegments(page)
	return len(items) == 0, err
}

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

type GetResult struct{ gophercloud.Result }

func (r GetResult) Extract() (*Segment, error) {
	var s struct {
		Item *Segment `json:"segment"`
	}
	err := r.ExtractInto(&s)
	return s.Item, err
}
