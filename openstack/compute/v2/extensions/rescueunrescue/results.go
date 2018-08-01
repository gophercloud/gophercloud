package rescueunrescue

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// RescueResult is the response from a Rescue operation. Call its Extract
// method to retrieve adminPass for a rescued server.
type RescueResult struct {
	commonResult
}

// Extract interprets any RescueResult as an AdminPass, if possible.
func (r RescueResult) Extract() (string, error) {
	var s struct {
		AdminPass string `json:"adminPass"`
	}
	err := r.ExtractInto(&s)
	return s.AdminPass, err
}
