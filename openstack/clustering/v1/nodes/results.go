package nodes

import (
	"github.com/gophercloud/gophercloud"
)

// DeleteResult is the result from a Delete operation. Call ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
