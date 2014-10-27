package cdncontainers

import (
	"net/http"

	"github.com/rackspace/gophercloud"
)

type headerResult struct {
	gophercloud.Result
}

func (hr headerResult) ExtractHeader() (http.Header, error) {
	return hr.Header, hr.Err
}

// EnableResult represents the result of a get operation.
type EnableResult struct {
	headerResult
}
