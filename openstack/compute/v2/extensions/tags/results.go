package tags

import "github.com/gophercloud/gophercloud"

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

type ListResult struct {
	tagResult
}
