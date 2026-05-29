package usages

// UsagesPre138 represents the total resource consumption for a project for
// microversions 1.9 through 1.37. In these versions, usages are not grouped
// by consumer type.
type UsagesPre138 struct {
	// Usages maps resource class names to the total integer amount consumed.
	Usages map[string]int `json:"usages"`
}

// ExtractPre138 interprets a GetResult as UsagesPre138 (microversions 1.9–1.37).
func (r GetResult) ExtractPre138() (*UsagesPre138, error) {
	var s UsagesPre138
	err := r.ExtractInto(&s)
	return &s, err
}
