package swauth

import (
	"github.com/gophercloud/gophercloud"
)

// GetAuthResult temporarily contains the response from a Swauth
// authentication call.
type GetAuthResult struct {
	gophercloud.HeaderResult
}

// GetAuthHeader contains the authentication information from a Swauth
// authentication request.
type GetAuthHeader struct {
	Token      string `json:"X-Auth-Token"`
	StorageURL string `json:"X-Storage-Url"`
	CDNURL     string `json:"X-CDN-Management-Url"`
}

// Extract is a method that attempts to interpret any Swauth authentication
// response as a GetAuthHeader struct.
func (r GetAuthResult) Extract() (*GetAuthHeader, error) {
	var s *GetAuthHeader
	err := r.ExtractInto(&s)
	return s, err
}
