package availabilityzones

import "github.com/gophercloud/gophercloud"

// AvailabilityZone contains all the information associated with an OpenStack
// AvailabilityZone.
type AvailabilityZone struct {
	// The availability zone ID.
	ID string `json:"id"`
	// The name of the availability zone.
	Name string `json:"name"`
	// The date and time stamp when the availability zone was created.
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	// The date and time stamp when the availability zone was updated.
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

// ListResult contains the response body and error from a List request.
type ListResult struct {
	commonResult
}

// Extract will get the AvailabilityZone objects out of the shareTypeAccessResult object.
func (r ListResult) Extract() ([]AvailabilityZone, error) {
	var a struct {
		AvailabilityZone []AvailabilityZone `json:"availability_zones"`
	}
	err := r.ExtractInto(&a)
	return a.AvailabilityZone, err
}
