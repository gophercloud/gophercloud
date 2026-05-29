package quotas

import (
	"github.com/gophercloud/gophercloud/v2"
)

// QuotaDetail represents a Quota API resource.
type QuotaDetail struct {
	Resource string
	Limit    int
	InUse    int `json:"in_use"`
	Reserved int
}

// GetResult is the result of a Get operation. Call its Extract method to
// interpret the result as a []QuotaDetail.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as a []QuotaDetail.
// An error is returned if the original call or the extraction failed.
func (r GetResult) Extract() ([]QuotaDetail, error) {
	var s struct {
		Quotas []QuotaDetail `json:"quotas"`
	}
	err := r.ExtractInto(&s)
	return s.Quotas, err
}
