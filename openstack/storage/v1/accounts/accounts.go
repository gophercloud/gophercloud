package accounts

import (
	"strings"
)

type UpdateOpts struct {
	Metadata map[string]string
	Headers  map[string]string
}

type GetOpts struct {
	Headers map[string]string
}

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
