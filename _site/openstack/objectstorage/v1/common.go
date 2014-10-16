package v1

import (
	"net/http"
)

// CommonResult is a structure that contains the response and error of a call to an
// object storage endpoint.
type CommonResult struct {
	Resp *http.Response
	Err  error
}

// ExtractHeaders will extract and return the headers from a *http.Response.
func (cr CommonResult) ExtractHeaders() (http.Header, error) {
	if cr.Err != nil {
		return nil, cr.Err
	}

	var headers http.Header
	if cr.Err != nil {
		return headers, cr.Err
	}
	return cr.Resp.Header, nil
}
