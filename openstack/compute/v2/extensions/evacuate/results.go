package evacuate

import (
	"github.com/gophercloud/gophercloud"
)

// EvacuateResult is the response from an Evacuate operation. Call its ExtractErr
// method to determine if the request suceeded or failed.
type EvacuateResult struct {
	gophercloud.ErrResult
}
