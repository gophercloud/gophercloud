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
