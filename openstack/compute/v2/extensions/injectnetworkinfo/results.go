package injectnetworkinfo

import (
	"github.com/gophercloud/gophercloud"
)

// InjectNetworkResult is the response of a InjectNetworkInfo operation. Call
// its ExtractErr method to determine if the request suceeded or failed.
type InjectNetworkResult struct {
	gophercloud.ErrResult
}
