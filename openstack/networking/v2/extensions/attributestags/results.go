package attributestags

import (
	"github.com/gophercloud/gophercloud"
)

type tagResult struct {
	gophercloud.Result
}

// Extract interprets tagResult to return the list of tags
func (r tagResult) Extract() ([]string, error) {
	var s struct {
		Tags []string `json:"tags"`
	}
	err := r.ExtractInto(&s)
	return s.Tags, err
}

// ReplaceAllResult represents the result of a replace operation.
// Call its Extract method to interpret it as a slice of strings.
type ReplaceAllResult struct {
	tagResult
}

type ListResult struct {
	tagResult
}

// DeleteResult is the result from a Delete/DeleteAll operation.
// Call its ExtractErr method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
