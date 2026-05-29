package segments

import (
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Segment model
type Segment struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	NetworkID       string    `json:"network_id"`
	NetworkType     string    `json:"network_type"`
	PhysicalNetwork string    `json:"physical_network"`
	SegmentationID  int       `json:"segmentation_id"`
	RevisionNumber  int       `json:"revision_number"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SegmentPage wraps a page of segments.
type SegmentPage struct {
	pagination.LinkedPageBase
}

func (r SegmentPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractSegments(r)
	return len(is) == 0, err
}

func ExtractSegments(r pagination.Page) ([]Segment, error) {
	var s []Segment
	err := ExtractSegmentsInto(r, &s)
	return s, err
}

// ExtractSegmentsInto extracts the elements into a slice of Segment structs.
func ExtractSegmentsInto(r pagination.Page, v any) error {
	return r.(SegmentPage).ExtractIntoSlicePtr(v, "segments")
}

// Segment results
type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Segment, error) {
	var s Segment
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "segment")
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
