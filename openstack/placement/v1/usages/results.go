package usages

import "github.com/gophercloud/gophercloud/v2"

// ConsumerTypeUsage maps resource class names to the total integer amount
// consumed by consumers of a given type. The special key "consumer_count"
// holds the number of consumers of this type.
type ConsumerTypeUsage map[string]int

// Usages represents the total resource consumption for a project, grouped by
// consumer type. This is the response shape for microversion 1.38 and later.
//
// Each key in the Usages map is a consumer type name (e.g. "INSTANCE").
// The value is a ConsumerTypeUsage containing resource class totals and a
// consumer_count entry.
type Usages struct {
	// Usages maps consumer type names to their aggregated resource usage.
	Usages map[string]ConsumerTypeUsage `json:"usages"`
}

// GetResult is the result of a Get operation. Call its Extract method
// to interpret it as Usages (microversion 1.38+), or ExtractPre138 for
// earlier microversions.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as Usages (microversion 1.38+).
func (r GetResult) Extract() (*Usages, error) {
	var s Usages
	err := r.ExtractInto(&s)
	return &s, err
}
