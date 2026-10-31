package networksegmentranges

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a network segment range result as a NetworkSegmentRange.
func (r commonResult) Extract() (*NetworkSegmentRange, error) {
	var s struct {
		NetworkSegmentRange *NetworkSegmentRange `json:"network_segment_range"`
	}
	err := r.ExtractInto(&s)
	return s.NetworkSegmentRange, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// NetworkSegmentRange represents a network segment range.
type NetworkSegmentRange struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Default         bool           `json:"default"`
	Shared          bool           `json:"shared"`
	ProjectID       string         `json:"project_id"`
	NetworkType     string         `json:"network_type"`
	PhysicalNetwork string         `json:"physical_network"`
	Minimum         int            `json:"minimum"`
	Maximum         int            `json:"maximum"`
	Used            map[int]string `json:"used"`
	Available       []int          `json:"available"`
}

// NetworkSegmentRangePage is the page returned by a pager when traversing over a collection of network segment ranges.
type NetworkSegmentRangePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a NetworkSegmentRangePage struct is empty.
func (r NetworkSegmentRangePage) IsEmpty() (bool, error) {
	is, err := ExtractNetworkSegmentRanges(r)
	return len(is) == 0, err
}

// ExtractNetworkSegmentRanges accepts a Page struct and extracts the elements into a slice of NetworkSegmentRange structs.
func ExtractNetworkSegmentRanges(r pagination.Page) ([]NetworkSegmentRange, error) {
	var s struct {
		NetworkSegmentRanges []NetworkSegmentRange `json:"network_segment_ranges"`
	}
	err := (r.(NetworkSegmentRangePage)).ExtractInto(&s)
	return s.NetworkSegmentRanges, err
}
