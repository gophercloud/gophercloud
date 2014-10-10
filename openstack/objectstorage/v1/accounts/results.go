package accounts

import (
	"net/http"
	"strings"
)

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func ExtractMetadata(gr GetResult) map[string]string {
	metadata := make(map[string]string)
	for k, v := range gr.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata
}
