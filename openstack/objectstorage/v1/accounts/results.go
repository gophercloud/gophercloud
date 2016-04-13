package accounts

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

// UpdateResult is returned from a call to the Update function.
type UpdateResult struct {
	gophercloud.HeaderResult
}

// UpdateHeader represents the headers returned in the response from an Update request.
type UpdateHeader struct {
	ContentLength string                  `json:"Content-Length"`
	ContentType   string                  `json:"Content-Type"`
	Date          gophercloud.JSONRFC1123 `json:"Date"`
	TransID       string                  `json:"X-Trans-Id"`
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (ur UpdateResult) Extract() (*UpdateHeader, error) {
	var uh *UpdateHeader
	err := ur.ExtractInto(&uh)
	return uh, err
}

// GetHeader represents the headers returned in the response from a Get request.
type GetHeader struct {
	BytesUsed      string                  `json:"X-Account-Bytes-Used"`
	ContainerCount string                  `json:"X-Account-Container-Count"`
	ContentLength  string                  `json:"Content-Length"`
	ContentType    string                  `json:"Content-Type"`
	Date           gophercloud.JSONRFC1123 `json:"Date"`
	ObjectCount    string                  `json:"X-Account-Object-Count"`
	TransID        string                  `json:"X-Trans-Id"`
	TempURLKey     string                  `json:"X-Account-Meta-Temp-URL-Key"`
	TempURLKey2    string                  `json:"X-Account-Meta-Temp-URL-Key-2"`
}

// GetResult is returned from a call to the Get function.
type GetResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (r GetResult) Extract() (*GetHeader, error) {
	var s *GetHeader
	err := r.ExtractInto(&s)
	return s, err
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (r GetResult) ExtractMetadata() (map[string]string, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	metadata := make(map[string]string)
	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}
