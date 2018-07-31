package rescueunrescue

import "github.com/gophercloud/gophercloud"

// RescueResult is the response from a Rescue operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type RescueResult struct {
	gophercloud.ErrResult
}

// Extract interprets any RescueResult as an AdminPass, if possible.
func (r RescueResult) Extract() (string, error) {
	var s struct {
		AdminPass string `json:"adminPass"`
	}
	err := r.ExtractInto(&s)
	return s.AdminPass, err
}
