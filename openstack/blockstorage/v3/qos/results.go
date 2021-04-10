package qos

import "github.com/gophercloud/gophercloud"

// QoS contains all the information associated with an OpenStack QoS specification.
type QoS struct {
	// Name is the name of the QoS.
	Name string `json:"name"`
	// Unique identifier for the QoS.
	ID string `json:"id"`
	// Consumer of QoS
	Consumer string `json:"consumer"`
	// Arbitrary key-value pairs defined by the user.
	Specs map[string]string `json:"specs"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the QoS object out of the commonResult object.
func (r commonResult) Extract() (*QoS, error) {
	var s QoS
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a QoS struct
func (r commonResult) ExtractInto(qos interface{}) error {
	return r.Result.ExtractIntoStructPtr(qos, "qos_specs")
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}
