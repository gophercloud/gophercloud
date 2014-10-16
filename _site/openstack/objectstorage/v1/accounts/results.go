package accounts

import (
	"strings"

	objectstorage "github.com/rackspace/gophercloud/openstack/objectstorage/v1"
)

// GetResult is returned from a call to the Get function. See v1.CommonResult.
type GetResult struct {
	objectstorage.CommonResult
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (gr GetResult) ExtractMetadata() (map[string]string, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}

	metadata := make(map[string]string)
	for k, v := range gr.Resp.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}

// UpdateResult is returned from a call to the Update function. See v1.CommonResult.
type UpdateResult struct {
	objectstorage.CommonResult
}
