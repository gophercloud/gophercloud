package volumeactions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// AttachResult contains the response body and error from a Get request.
type AttachResult struct {
	commonResult
}

// DetachResult contains the response body and error from a Get request.
type DetachResult struct {
	commonResult
}

// ReserveResult contains the response body and error from a Get request.
type ReserveResult struct {
	commonResult
}

// UnreserveResult contains the response body and error from a Get request.
type UnreserveResult struct {
	commonResult
}

// InitializeConnectionResult contains the response body and error from a Get request.
type InitializeConnectionResult struct {
	commonResult
}

// TerminateConnectionResult contains the response body and error from a Get request.
type TerminateConnectionResult struct {
	commonResult
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (map[string]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res map[string]interface{}

	err := mapstructure.Decode(r.Body, &res)

	return res, err
}
