package accounts

import (
	"net/http"
	"strings"

	"github.com/rackspace/gophercloud"
)

type headerResult struct {
	gophercloud.Result
}

func (hr headerResult) ExtractHeader() (http.Header, error) {
	return hr.Header, hr.Err
}

// GetResult is returned from a call to the Get function.
type GetResult struct {
	headerResult
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (gr GetResult) ExtractMetadata() (map[string]string, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}

	metadata := make(map[string]string)
	for k, v := range gr.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}

// UpdateResult is returned from a call to the Update function.
type UpdateResult struct {
	headerResult
}
