package lockunlock

import (
	"github.com/gophercloud/gophercloud"
)

// LockResult and UnlockResult are the responses from a Lock and Unlock operations respectively. Call its ExtractErr
// method to determine if the suceeded or failed.
type LockResult struct {
	gophercloud.ErrResult
}

type UnlockResult struct {
	gophercloud.ErrResult
}
