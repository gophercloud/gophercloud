package accounts

import (
	"net/http"
	"strings"

	"github.com/rackspace/gophercloud"
)

type commonResult struct {
	gophercloud.CommonResult
	Resp *http.Response
}

// GetResult represents the result of a create operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (res commonResult) ExtractMetadata() map[string]string {
	metadata := make(map[string]string)

	for k, v := range res.Resp.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}

	return metadata
}
